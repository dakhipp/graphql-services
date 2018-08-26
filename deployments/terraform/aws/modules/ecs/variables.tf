variable "environment" {
  description = "The environment"
  default     = "staging"
}

variable "vpc_id" {
  description = "The VPC id"
}

variable "availability_zones" {
  type        = "list"
  description = "The azs to use"
}

variable "security_groups_ids" {
  type        = "list"
  description = "The SGs to use"
}

variable "subnets_ids" {
  type        = "list"
  description = "The private subnets to use"
}

variable "public_subnet_ids" {
  type        = "list"
  description = "The private subnets to use"
}

variable "graphql_port" {
  description = "The port the graphql container will start on"
}

variable "auth_port" {
  description = "The port the auth container will start on"
}

variable "playground_enabled" {
  description = "Enable GraphQL playground"
}

variable "psql_addr" {
  description = "The database endpoint"
}

variable "psql_user" {
  description = "The database username"
}

variable "psql_pass" {
  description = "The database password"
}

variable "psql_db" {
  description = "The database that the app will use"
}

variable "psql_ssl" {
  description = "Database ssl enabled or disabled"
}

variable "psql_port" {
  description = "The database port"
}

variable "graphql_repository_name" {
  description = "The name of the graphql ECR repisitory"
}

variable "auth_repository_name" {
  description = "The name of the graphql ECR repisitory"
}
