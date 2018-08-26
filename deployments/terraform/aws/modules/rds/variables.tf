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

variable "psql_db" {
  description = "The database name"
}

variable "psql_user" {
  description = "The username of the database"
}

variable "psql_pass" {
  description = "The password of the database"
}

variable "subnet_ids" {
  type        = "list"
  description = "Subnet ids"
}

variable "vpc_id" {
  description = "The VPC id"
}
