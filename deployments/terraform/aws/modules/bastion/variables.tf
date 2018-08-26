variable "environment" {
  description = "The environment"
  default     = "staging"
}

variable "bastion_key_name" {
  description = "A name for the SSH key"
  default     = "bastion-key-staging"
}

variable "bastion_public_key" {
  description = "Public SSH key supplied by terraform.tfvars file, orginates from terraform.tfvars file"
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
