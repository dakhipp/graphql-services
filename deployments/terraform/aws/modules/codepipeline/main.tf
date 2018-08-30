locals {
  named_codebuild = "graphql_codebuild_${var.environment}"
}

// S3 bucket holding source downloaded from github and final artifact manifest for codeDeploy
resource "aws_s3_bucket" "source" {
  bucket        = "${var.artifact_bucket_name}"
  acl           = "private"
  force_destroy = true
}

// CodePipeline IAM role
resource "aws_iam_role" "codepipeline_role" {
  name = "codepipeline_role_${var.environment}"

  assume_role_policy = "${file("${path.module}/policies/codepipeline_role.json")}"
}

// CodePipeline policy with templating variables
data "template_file" "codepipeline_policy" {
  template = "${file("${path.module}/policies/codepipeline.json")}"

  vars {
    aws_s3_bucket_arn = "${aws_s3_bucket.source.arn}"
  }
}

// Rendered CodePipeline policy
resource "aws_iam_role_policy" "codepipeline_policy" {
  name   = "codepipeline_policy"
  role   = "${aws_iam_role.codepipeline_role.id}"
  policy = "${data.template_file.codepipeline_policy.rendered}"
}

/*
/* CodeBuild
*/
// CodeBuild IAM role
resource "aws_iam_role" "codebuild_role" {
  name               = "codebuild_role_${var.environment}"
  assume_role_policy = "${file("${path.module}/policies/codebuild_role.json")}"
}

// CodeBuild policy with templating variables
data "template_file" "codebuild_policy" {
  template = "${file("${path.module}/policies/codebuild_policy.json")}"

  vars {
    aws_s3_bucket_arn = "${aws_s3_bucket.source.arn}"
  }
}

// Rendered CodeBuild policy template
resource "aws_iam_role_policy" "codebuild_policy" {
  name   = "codebuild_policy"
  role   = "${aws_iam_role.codebuild_role.id}"
  policy = "${data.template_file.codebuild_policy.rendered}"
}

// CodeBuild buildspec with templating variables
data "template_file" "buildspec" {
  template = "${file("${path.module}/buildspec.yml")}"

  vars {
    region             = "${var.region}"
    environment        = "${var.environment}"
    graphql_name       = "${var.graphql_service_name}"
    graphql_repo_url   = "${var.graphql_repo_url}"
    auth_name          = "${var.auth_service_name}"
    auth_repo_url      = "${var.auth_repo_url}"
    migration_name     = "${var.migration_service_name}"
    migration_repo_url = "${var.migration_repo_url}"
    cluster_name       = "${var.ecs_cluster_name}"
    subnet_id          = "${var.run_task_subnet_id}"
    security_group_ids = "${join(",", var.run_task_security_group_ids)}"
  }
}

// CodeBuild configuration
resource "aws_codebuild_project" "graphql_build" {
  name          = "${local.named_codebuild}"
  build_timeout = "10"
  service_role  = "${aws_iam_role.codebuild_role.arn}"

  artifacts {
    type = "CODEPIPELINE"
  }

  environment {
    compute_type = "BUILD_GENERAL1_SMALL"

    // https://docs.aws.amazon.com/codebuild/latest/userguide/build-env-ref-available.html
    image           = "aws/codebuild/golang:1.10"
    type            = "LINUX_CONTAINER"
    privileged_mode = true
  }

  source {
    type      = "CODEPIPELINE"
    buildspec = "${data.template_file.buildspec.rendered}"
  }
}

// CodePipeline configuration
resource "aws_codepipeline" "pipeline" {
  name     = "${local.named_codebuild}"
  role_arn = "${aws_iam_role.codepipeline_role.arn}"

  artifact_store {
    location = "${aws_s3_bucket.source.bucket}"
    type     = "S3"
  }

  stage {
    name = "Source"

    action {
      name             = "Source"
      category         = "Source"
      owner            = "ThirdParty"
      provider         = "GitHub"
      version          = "1"
      output_artifacts = ["source"]

      configuration {
        Owner      = "${var.github_user}"
        Repo       = "${var.github_repo}"
        Branch     = "${var.github_branch}"
        OAuthToken = "${var.github_oauth}"
      }
    }
  }

  stage {
    name = "Build"

    action {
      name             = "Build"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      version          = "1"
      input_artifacts  = ["source"]
      output_artifacts = ["imagedefinitions"]

      configuration {
        ProjectName = "graphql_codebuild_${var.environment}"
      }
    }
  }

  stage {
    name = "Production"

    action {
      name            = "Deploy"
      category        = "Deploy"
      owner           = "AWS"
      provider        = "ECS"
      input_artifacts = ["imagedefinitions"]
      version         = "1"

      configuration {
        ClusterName = "${var.ecs_cluster_name}"
        ServiceName = "${var.ecs_service_name}"
        FileName    = "imagedefinitions.json"
      }
    }
  }
}
