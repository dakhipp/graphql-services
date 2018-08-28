variable "domain" {
  description = "The domain to create an a record on"
}

variable "alb_dns_name" {
  description = "The environment"
  default     = "The ALB DNS name provided from the ECS module"
}

variable "alb_zone_id" {
  description = "The environment"
  default     = "The ALB zone ID provided from the ECS module"
}

variable "route53_zone_id" {
  description = "Zone ID for an existing Route53 hosted zone, originates from terraform.tfvars file"
}
