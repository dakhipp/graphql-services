output "graphql_repo_url" {
  value = "${aws_ecr_repository.graphql_repo.repository_url}"
}

output "auth_repo_url" {
  value = "${aws_ecr_repository.auth_repo.repository_url}"
}

output "migration_repo_url" {
  value = "${aws_ecr_repository.migration_repo.repository_url}"
}

output "cluster_name" {
  value = "${aws_ecs_cluster.cluster.name}"
}

output "service_name" {
  value = "${aws_ecs_service.graphql_web.name}"
}

output "alb_dns_name" {
  value = "${aws_alb.alb_graphql.dns_name}"
}

output "alb_zone_id" {
  value = "${aws_alb.alb_graphql.zone_id}"
}

output "security_group_id" {
  value = "${aws_security_group.ecs_service.id}"
}
