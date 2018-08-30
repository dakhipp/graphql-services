variable "environment" {
  description = "The environment"
  default     = "staging"
}

variable "bastion_key_name" {
  description = "A name for the SSH key"
  default     = "bastion-key-staging"
}

variable "psql_root_db" {
  description = "Database name supplied by terraform.tfvars file"
}

variable "psql_root_user" {
  description = "The root database user supplied by terraform.tfvars file"
}

variable "psql_root_pass" {
  description = "The root database password supplied by terraform.tfvars file"
}

variable "psql_web_db" {
  description = "The database name for this app supplied by terraform.tfvars file"
}

variable "psql_web_user" {
  description = "A limited database user for this app supplied by terraform.tfvars file"
}

variable "psql_web_pass" {
  description = "A database password for the limited database user for this app, supplied by terraform.tfvars file"
}

variable "psql_port" {
  description = "PostgreSQL port supplied by terraform.tfvars file"
}

variable "bastion_public_key" {
  description = "Public SSH key supplied by terraform.tfvars file, orginates from terraform.tfvars file"
}

variable "psql_addr" {
  description = "The database endpoint, originates from RDS module"
}

variable "subnet_id" {
  description = "The ID of the subnet the server should be placed in, originates from VPC module"
}

variable "vpc_id" {
  description = "The ID of the VPC the server should be placed in, originates from VPC module"
}

variable "rds_sg" {
  description = "Security group that allows access to RDS, originates from RDS module"
}
