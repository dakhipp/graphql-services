variable "environment" {
  description = "The environment"
  default     = "staging"
}

variable "region" {
  description = "The region to use"
  default     = "us-west-2"
}

variable "graphql_service_name" {
  description = "The name of the GraphQL ECR repository"
  default     = "graphql"
}

variable "auth_service_name" {
  description = "The name of the auth ECR repository"
  default     = "auth"
}

variable "migration_service_name" {
  description = "The name of the migrations ECR repository"
  default     = "migration"
}

variable "github_oauth" {
  description = "Github OAuth token, originates from the tarraform.tfvars file"
}

variable "github_user" {
  description = "Github user, originates from the tarraform.tfvars file"
}

variable "github_repo" {
  description = "Github repo, originates from the tarraform.tfvars file"
}

variable "github_branch" {
  description = "Github branch, originates from the tarraform.tfvars file"
}

variable "artifact_bucket_name" {
  description = "The name of the bucket to store codepipeline artifacts in, originates from the tarraform.tfvars file"
}

variable "graphql_repo_url" {
  description = "The url of the GraphQL ECR repository, originates from the ECS module"
}

variable "auth_repo_url" {
  description = "The url of the auth ECR repository, originates from the ECS module"
}

variable "migration_repo_url" {
  description = "The url of the migraions ECR repository, originates from the ECS module"
}

variable "ecs_service_name" {
  description = "The ECS service that will be deployed, originates from the ECS module"
}

variable "ecs_cluster_name" {
  description = "The cluster that we will deploy, originates from the ECS module"
}

variable "run_task_subnet_id" {
  description = "The subnet Id where single run task will be executed, originates from the VPC module"
}

variable "run_task_security_group_ids" {
  type        = "list"
  description = "The security group Ids attached where the single run task will be executed, originates from RDS, VPC, and ECS modules"
}
