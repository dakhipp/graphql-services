variable "bastion_public_key" {
  description = "Public SSH key supplied by terraform.tfvars file"
}

variable "psql_user" {
  description = "Database user supplied by terraform.tfvars file"
}

variable "psql_pass" {
  description = "Database password supplied by terraform.tfvars file"
}

variable "psql_db" {
  description = "Database name supplied by terraform.tfvars file"
}

variable "psql_ssl" {
  description = "Database ssl enabled or disabled supplied by terraform.tfvars file"
}

variable "psql_port" {
  description = "PostgreSQL port supplied by terraform.tfvars file"
}

variable "graphql_port" {
  description = "Port the GraphQL service container will start on supplied by terraform.tfvars file"
}

variable "auth_port" {
  description = "Port the auth service container will start on supplied by terraform.tfvars file"
}

variable "playground_enabled" {
  description = "Whether or not the GraphQL playground should be enabled, supplied by terraform.tfvars file"
}
