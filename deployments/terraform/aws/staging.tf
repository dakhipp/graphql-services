locals {
  region             = "us-west-2"
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]

  // TODO: don't do this.
  db_name = "psql"
  db_user = "root"
  db_pass = "toor1234"
}

/*====
Provider
======*/
provider "aws" {
  region = "${local.region}"
}

/*====
Remote State Config
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
  environment          = "staging"
  vpc_cidr             = "10.0.0.0/16"
  public_subnets_cidr  = ["10.0.0.0/24", "10.0.1.0/24", "10.0.2.0/24"]
  private_subnets_cidr = ["10.0.100.0/24", "10.0.101.0/24", "10.0.102.0/24"]
  region               = "${local.region}"
  availability_zones   = "${local.availability_zones}"
  key_name             = "production_key"
}

module "rds" {
  source                  = "./modules/rds"
  environment             = "staging"
  allocated_storage       = "20"
  instance_class          = "db.t2.micro"
  engine                  = "postgres"
  engine_version          = "10.4"
  backup_retention_period = "7"
  backup_window           = "02:00-03:00"
  multi_az                = false
  apply_immediately       = true
  database_name           = "${local.db_name}"
  database_username       = "${local.db_user}"
  database_password       = "${local.db_pass}"
  subnet_ids              = ["${module.vpc.private_subnets_id}"]
  vpc_id                  = "${module.vpc.vpc_id}"
}
