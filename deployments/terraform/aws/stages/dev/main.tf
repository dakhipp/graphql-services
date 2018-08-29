locals {
  region = "us-west-2"
}

/*====
Provider
======*/
provider "aws" {
  region  = "${local.region}"
  version = "1.33"
}

/*====
Remote State Config, S3 bucket must exist before running, values cannot be parameterized
======*/
terraform {
  backend "s3" {
    bucket = "graphql-service-state"
    region = "us-west-2"
    key    = "terraform-state-dev"   // MUST be unique betweens stages

    # Slowing down builds, don't need this for now
    # dynamodb_table = "graphql-lock-table"   // MUST be unique betweens stages
  }
}

module "vpc" {
  source               = "../../modules/vpc"
  region               = "${local.region}"
  environment          = "${var.environment}"
  vpc_cidr             = "10.0.0.0/16"
  public_subnets_cidr  = ["10.0.0.0/24", "10.0.1.0/24", "10.0.2.0/24"]
  private_subnets_cidr = ["10.0.100.0/24", "10.0.101.0/24", "10.0.102.0/24"]
  availability_zones   = ["us-west-2a", "us-west-2b", "us-west-2c"]
}

module "rds" {
  source                  = "../../modules/rds"
  environment             = "${var.environment}"
  allocated_storage       = "20"
  instance_class          = "db.t2.micro"
  engine                  = "postgres"
  engine_version          = "10.4"
  backup_retention_period = "7"
  backup_window           = "02:00-03:00"
  multi_az                = false
  apply_immediately       = true
  psql_db                 = "${var.psql_db}"
  psql_user               = "${var.psql_user}"
  psql_pass               = "${var.psql_pass}"
  psql_port               = "${var.psql_port}"
  subnet_ids              = ["${module.vpc.private_subnets_id}"]
  vpc_id                  = "${module.vpc.vpc_id}"
}

module "bastion" {
  source             = "../../modules/bastion"
  environment        = "${var.environment}"
  bastion_key_name   = "bastion-key-${var.environment}"
  bastion_public_key = "${var.bastion_public_key}"
  subnet_id          = "${module.vpc.public_subnets_id[0]}"
  vpc_id             = "${module.vpc.vpc_id}"
  rds_sg             = "${module.rds.db_access_sg_id}"
}

module "ecs" {
  source                 = "../../modules/ecs"
  environment            = "${var.environment}"
  domain                 = "${var.domain}"
  ssl_identifier         = "${var.ssl_identifier}"
  availability_zones     = "${var.graphql_service_name}"
  graphql_service_name   = "${var.graphql_service_name}"
  auth_service_name      = "${var.auth_service_name}"
  migration_service_name = "${var.migration_service_name}"
  vpc_id                 = "${module.vpc.vpc_id}"
  subnets_ids            = ["${module.vpc.private_subnets_id}"]
  public_subnet_ids      = ["${module.vpc.public_subnets_id}"]

  security_groups_ids = [
    "${module.vpc.security_groups_ids}",
    "${module.rds.db_access_sg_id}",
  ]

  // graphql env vars
  graphql_port       = "${var.graphql_port}"
  playground_enabled = "${var.playground_enabled}"

  // auth env vars
  auth_port = "${var.auth_port}"
  psql_user = "${var.psql_user}"
  psql_pass = "${var.psql_pass}"
  psql_db   = "${var.psql_db}"
  psql_ssl  = "${var.psql_ssl}"
  psql_port = "${var.psql_port}"
  psql_addr = "${module.rds.rds_address}"
}

module "codepipeline" {
  source                      = "../../modules/codepipeline"
  region                      = "${local.region}"
  environment                 = "${var.environment}"
  github_oauth                = "${var.github_oauth}"
  github_user                 = "${var.github_user}"
  github_repo                 = "${var.github_repo}"
  github_branch               = "${var.github_branch}"
  artifact_bucket_name        = "${var.artifact_bucket_name}"
  graphql_service_name        = "${var.graphql_service_name}"
  auth_service_name           = "${var.auth_service_name}"
  migration_service_name      = "${var.migration_service_name}"
  graphql_repo_url            = "${module.ecs.graphql_repo_url}"
  auth_repo_url               = "${module.ecs.auth_repo_url}"
  migration_repo_url          = "${module.ecs.migration_repo_url}"
  ecs_service_name            = "${module.ecs.service_name}"
  ecs_cluster_name            = "${module.ecs.cluster_name}"
  run_task_subnet_id          = "${module.vpc.private_subnets_id[0]}"
  run_task_security_group_ids = ["${module.rds.db_access_sg_id}", "${module.vpc.security_groups_ids}", "${module.ecs.security_group_id}"]
}

module "route53" {
  source          = "../../modules/route53"
  domain          = "${var.domain}"
  route53_zone_id = "${var.route53_zone_id}"
  alb_dns_name    = "${module.ecs.alb_dns_name}"
  alb_zone_id     = "${module.ecs.alb_zone_id}"
}