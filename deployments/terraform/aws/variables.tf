variable "bastion_public_key" {
  description = "Public SSH key supplied by terraform.tfvars file"
}

variable "db_name" {
  description = "Database name supplied by terraform.tfvars file"
}

variable "db_user" {
  description = "Database user supplied by terraform.tfvars file"
}

variable "db_pass" {
  description = "Database password supplied by terraform.tfvars file"
}
