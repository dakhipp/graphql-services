variable "bastion_public_key" {
  description = "Public SSH key supplied by terraform.tfvars file"
}

variable "github_oauth" {
  description = "Github OAuth token supplied by terraform.tfvars file"
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

variable "domain" {
  description = "The domain to create an a record on, supplied by terraform.tfvars file"
}

variable "ssl_identifier" {
  description = "A domain added to an ACM certificate, sometimes the same as the domain variable, but if using a wildcard cert it might not be. Originates from terraform.tfvars file"
}

variable "route53_zone_id" {
  description = "Zone ID for an existing Route53 hosted zone, originates from terraform.tfvars file"
}
