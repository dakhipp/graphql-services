locals {
  region             = "us-west-2"
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
  environment        = "staging"
}

/*====
Provider
======*/
provider "aws" {
  region  = "${local.region}"
  version = "1.33"
}

/*====
Remote State Config, S3 bucket must exist before running
======*/
terraform {
  backend "s3" {
    bucket = "terraform-remote-state-123"
    region = "us-west-2"
    key    = "staging"
  }
}

module "vpc" {
  source               = "./modules/vpc"
  environment          = "${local.environment}"
  vpc_cidr             = "10.0.0.0/16"
  public_subnets_cidr  = ["10.0.0.0/24", "10.0.1.0/24", "10.0.2.0/24"]
  private_subnets_cidr = ["10.0.100.0/24", "10.0.101.0/24", "10.0.102.0/24"]
  region               = "${local.region}"
  availability_zones   = "${local.availability_zones}"
  key_name             = "production_key"
}

module "rds" {
  source                  = "./modules/rds"
  environment             = "${local.environment}"
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
  source             = "./modules/bastion"
  environment        = "${local.environment}"
  bastion_key_name   = "bastion-key-${local.environment}"
  bastion_public_key = "${var.bastion_public_key}"
  subnet_id          = "${module.vpc.public_subnets_id[0]}"
  vpc_id             = "${module.vpc.vpc_id}"
  rds_sg             = "${module.rds.db_access_sg_id}"
}

module "ecs" {
  source                    = "./modules/ecs"
  environment               = "${local.environment}"
  availability_zones        = "${local.production_availability_zones}"
  graphql_repository_name   = "${local.environment}/graphql"
  auth_repository_name      = "${local.environment}/auth"
  migration_repository_name = "${local.environment}/migrations"
  vpc_id                    = "${module.vpc.vpc_id}"
  subnets_ids               = ["${module.vpc.private_subnets_id}"]
  public_subnet_ids         = ["${module.vpc.public_subnets_id}"]

  security_groups_ids = [
    "${module.vpc.security_groups_ids}",
    "${module.rds.db_access_sg_id}",
  ]

  // graphql env vars
  graphql_port       = "${var.graphql_port}"
  playground_enabled = "${var.playground_enabled}"

  // auth env vars
  auth_port = "${var.auth_port}"
  psql_addr = "${module.rds.rds_address}"
  psql_user = "${var.psql_user}"
  psql_pass = "${var.psql_pass}"
  psql_db   = "${var.psql_db}"
  psql_ssl  = "${var.psql_ssl}"
  psql_port = "${var.psql_port}"
}

module "codepipeline" {
  source                      = "./modules/codepipeline"
  environment                 = "${local.environment}"
  region                      = "${local.region}"
  graphql_repository_url      = "${module.ecs.graphql_repository_url}"
  auth_repository_url         = "${module.ecs.auth_repository_url}"
  migration_repository_url    = "${module.ecs.migration_repository_url}"
  ecs_service_name            = "${module.ecs.service_name}"
  ecs_cluster_name            = "${module.ecs.cluster_name}"
  run_task_subnet_id          = "${module.vpc.private_subnets_id[0]}"
  run_task_security_group_ids = ["${module.rds.db_access_sg_id}", "${module.vpc.security_groups_ids}", "${module.ecs.security_group_id}"]
}
