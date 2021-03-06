variable "environment" {
  description = "The environment"
  default     = "staging"
}

variable "allocated_storage" {
  description = "The storage size in GB"
  default     = "20"
}

variable "instance_class" {
  description = "The instance type"
  default     = "db.t2.micro"
}

variable "engine" {
  description = "The RDS service to start"
  default     = "postgres"
}

variable "engine_version" {
  description = "Version of the selected engine"
  default     = "10.4"
}

variable "backup_retention_period" {
  description = "Number of days to retain database backups"
  default     = "7"
}

variable "backup_window" {
  description = "UTC time period for backups to be made"
  default     = "02:00-03:00"
}

variable "multi_az" {
  description = "Muti-az allowed?"
  default     = false
}

variable "apply_immediately" {
  description = "Apply updates immediately?"
  default     = false
}

variable "psql_root_db" {
  description = "The root database name, originates from terraform.tfvars file"
}

variable "psql_root_user" {
  description = "The root username of the database, originates from terraform.tfvars file"
}

variable "psql_root_pass" {
  description = "The root password for the root user of the database, originates from terraform.tfvars file"
}

variable "psql_port" {
  description = "The port for the database, originates from terraform.tfvars file"
}

variable "subnet_ids" {
  type        = "list"
  description = "Subnet ids, originates from VPC module"
}

variable "vpc_id" {
  description = "The VPC ID, originates from VPC module"
}
