variable "environment" {
  description = "The environment"
  default     = "staging"
}

variable "graphql_service_name" {
  description = "The name of the GraphQL ECR repository"
}

variable "auth_service_name" {
  description = "The name of the auth ECR repository"
}

variable "migration_service_name" {
  description = "The name of the migrations ECR repository"
}

variable "github_oauth" {
  description = "Github OAuth token"
}

variable "github_user" {
  description = "Github user"
}

variable "github_repo" {
  description = "Github repo"
}

variable "github_branch" {
  description = "Github branch"
}

variable "region" {
  description = "The region to use"
}

variable "graphql_repo_url" {
  description = "The url of the GraphQL ECR repository"
}

variable "auth_repo_url" {
  description = "The url of the auth ECR repository"
}

variable "migration_repo_url" {
  description = "The url of the migraions ECR repository"
}

variable "ecs_service_name" {
  description = "The ECS service that will be deployed"
}

variable "ecs_cluster_name" {
  description = "The cluster that we will deploy"
}

variable "run_task_subnet_id" {
  description = "The subnet Id where single run task will be executed"
}

variable "run_task_security_group_ids" {
  type        = "list"
  description = "The security group Ids attached where the single run task will be executed"
}

variable "artifact_bucket_name" {
  description = "The name of the bucket to store codepipeline artifacts in"
}
