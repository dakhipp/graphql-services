variable "environment" {
  description = "The environment"
  default     = "staging"
}

variable "availability_zones" {
  type        = "list"
  description = "The azs to use"
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

variable "vpc_id" {
  description = "The VPC id, originates from the VPC service"
}

variable "subnets_ids" {
  type        = "list"
  description = "The private subnets to place our cluster resources in, originates from VPC service"
}

variable "public_subnet_ids" {
  type        = "list"
  description = "The public subnets to place our ALB in, originates from VPC service"
}

variable "security_groups_ids" {
  type        = "list"
  description = "The SGs to use, originates from VPC service and RDS service"
}

variable "graphql_port" {
  description = "The port the GraphQL container will start on, originates from terraform.tfvars file"
}

variable "playground_enabled" {
  description = "Enable GraphQL playground, originates from terraform.tfvars file"
}

variable "auth_port" {
  description = "The port the auth container will start on, originates from terraform.tfvars file"
}

variable "psql_url" {
  description = "The database url, originates from RDS service"
}

variable "psql_web_db" {
  description = "The database that the app will use, originates from terraform.tfvars file"
}

variable "psql_web_user" {
  description = "The database username, originates from terraform.tfvars file"
}

variable "psql_web_pass" {
  description = "The database password, originates from terraform.tfvars file"
}

variable "psql_addr" {
  description = "The database endpoint, originates from RDS service"
}

variable "domain" {
  description = "The domain to create an a record on, originates from terraform.tfvars file"
}

variable "ssl_identifier" {
  description = "A domain added to an ACM certificate, sometimes the same as the domain variable, but if using a wildcard cert it might not be. Originates from terraform.tfvars file"
}
