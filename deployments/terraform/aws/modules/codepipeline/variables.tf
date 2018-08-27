variable "environment" {
  description = "The environment"
  default     = "staging"
}

variable "github_oauth" {
  description = "Github OAuth token"
}

variable "region" {
  description = "The region to use"
}

variable "graphql_repository_url" {
  description = "The url of the GraphQL ECR repository"
}

variable "auth_repository_url" {
  description = "The url of the auth ECR repository"
}

variable "migration_repository_url" {
  description = "The url of the migraions ECR repository"
}

variable "ecs_service_name" {
  description = "The ECS service that will be deployed"
}

variable "ecs_cluster_name" {
  description = "The cluster that we will deploy"
}

variable "run_task_subnet_id" {
  description = "The subnet Id where single run task will be executed"
}

variable "run_task_security_group_ids" {
  type        = "list"
  description = "The security group Ids attached where the single run task will be executed"
}
