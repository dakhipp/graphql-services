/*====
Cloudwatch Log Group
======*/
resource "aws_cloudwatch_log_group" "graphql-services" {
  name = "graphql-services"

  tags {
    Environment = "${var.environment}"
    Application = "graphql-services"
  }
}

/*====
ECR repository to store our Docker images
======*/
// GraphQL docker image
resource "aws_ecr_repository" "graphql_repo" {
  name = "${var.graphql_repository_name}"
}

// Auth docker image
resource "aws_ecr_repository" "auth_repo" {
  name = "${var.auth_repository_name}"
}

// Migrations docker image
resource "aws_ecr_repository" "migration_repo" {
  name = "${var.migration_repository_name}"
}

/*====
ECS cluster
======*/
resource "aws_ecs_cluster" "cluster" {
  name = "${var.environment}-ecs-cluster"
}

/*====
ECS task definition
======*/

// The task definition template for the GraphQL service
data "template_file" "graphql_task" {
  template = "${file("${path.module}/tasks/graphql_task_definition.json")}"

  // Variables passed into the task definition template file
  vars {
    graphql_name    = "${var.graphql_repository_name}"
    graphql_image   = "${aws_ecr_repository.graphql_repo.repository_url}"
    auth_name       = "${var.auth_repository_name}"
    auth_image      = "${aws_ecr_repository.auth_repo.repository_url}"
    migration_name  = "${var.migration_repository_name}"
    migration_image = "${aws_ecr_repository.migration_repo.repository_url}"
    log_group       = "${aws_cloudwatch_log_group.graphql-services.name}"

    // GraphQL env vars
    graphql_port       = "${var.graphql_port}"
    playground_enabled = "${var.playground_enabled}"

    // Auth env vars
    auth_port = "${var.auth_port}"
    psql_addr = "${var.psql_addr}:${var.psql_port}"
    psql_user = "${var.psql_user}"
    psql_pass = "${var.psql_pass}"
    psql_db   = "${var.psql_db}"
    psql_ssl  = "${var.psql_ssl}"
  }
}

// Rendered and fully configured task definition for the GraphQL service
resource "aws_ecs_task_definition" "graphql_web" {
  depends_on               = ["aws_ecs_task_definition.graphql_web"]
  family                   = "${var.environment}_graphql"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = "${aws_iam_role.ecs_execution_role.arn}"
  task_role_arn            = "${aws_iam_role.ecs_execution_role.arn}"

  // FIXME: The use of 'replace' here is because of a bug in 'jsonencode'. turns ints to integers unless prepended with 'string:':
  //        https://github.com/hashicorp/terraform/issues/17033
  container_definitions = "${replace(replace(data.template_file.graphql_task.rendered, "/\"([0-9]+\\.?[0-9]*)\"/", "$1"), "string:", "")}"
}

/*====
App Load Balancer
======*/
resource "random_id" "target_group_sufix" {
  byte_length = 2
}

/* security group for ALB */
resource "aws_security_group" "graphql_inbound_sg" {
  name        = "${var.environment}-web-inbound-sg"
  description = "Allow HTTP from Anywhere into ALB"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  // TODO: look into this block
  ingress {
    from_port   = 8
    to_port     = 0
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "${var.environment}-web-inbound-sg"
  }
}

// SSL cert for HTTPS fetched by domain name
data "aws_acm_certificate" "ssl" {
  domain      = "${var.ssl_identifier}"
  types       = ["AMAZON_ISSUED"]
  most_recent = true
}

resource "aws_alb" "alb_graphql" {
  name            = "${var.environment}-alb-graphql-services"
  subnets         = ["${var.public_subnet_ids}"]
  security_groups = ["${var.security_groups_ids}", "${aws_security_group.graphql_inbound_sg.id}"]

  tags {
    Name        = "${var.environment}-alb-graphql-services"
    Environment = "${var.environment}"
  }
}

// Load balancer target group forwarding traffic on port 80 to
resource "aws_alb_target_group" "alb_target_group" {
  name        = "${var.environment}-alb-target-group-${random_id.target_group_sufix.hex}"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = "${var.vpc_id}"
  target_type = "ip"

  health_check {
    path = "/h"
  }

  lifecycle {
    create_before_destroy = true
  }
}

// HTTP listener on load balancer
resource "aws_alb_listener" "graphql-services-http" {
  load_balancer_arn = "${aws_alb.alb_graphql.arn}"
  port              = "80"
  protocol          = "HTTP"
  depends_on        = ["aws_alb_target_group.alb_target_group"]

  default_action {
    target_group_arn = "${aws_alb_target_group.alb_target_group.arn}"
    type             = "forward"
  }
}

// HTTPS listener on load balancer
resource "aws_alb_listener" "graphql-services-https" {
  load_balancer_arn = "${aws_alb.alb_graphql.arn}"
  port              = "443"
  protocol          = "HTTPS"
  depends_on        = ["aws_alb_target_group.alb_target_group"]
  ssl_policy        = "ELBSecurityPolicy-2015-05"
  certificate_arn   = "${data.aws_acm_certificate.ssl.arn}"

  default_action {
    type             = "forward"
    target_group_arn = "${aws_alb_target_group.alb_target_group.arn}"
  }
}

// Redirect all traffic to https
resource "aws_lb_listener_rule" "redirect_http_to_https" {
  listener_arn = "${aws_alb_listener.graphql-services-http.arn}"

  action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }

  condition {
    field  = "host-header"
    values = ["${var.domain}"]
  }
}

// Redirect load balancer URL
resource "aws_lb_listener_rule" "redirect_balancer_url" {
  listener_arn = "${aws_alb_listener.graphql-services-http.arn}"

  action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
      host        = "${var.domain}"
    }
  }

  condition {
    field  = "host-header"
    values = ["${aws_alb.alb_graphql.dns_name}"]
  }
}

/*
* IAM service role
*/
data "aws_iam_policy_document" "ecs_service_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs_role" {
  name               = "ecs_role"
  assume_role_policy = "${data.aws_iam_policy_document.ecs_service_role.json}"
}

data "aws_iam_policy_document" "ecs_service_policy" {
  statement {
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "elasticloadbalancing:Describe*",
      "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
      "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
      "ec2:Describe*",
      "ec2:AuthorizeSecurityGroupIngress",
    ]
  }
}

/* ecs service scheduler role */
resource "aws_iam_role_policy" "ecs_service_role_policy" {
  name = "ecs_service_role_policy"

  policy = "${data.aws_iam_policy_document.ecs_service_policy.json}"
  role   = "${aws_iam_role.ecs_role.id}"
}

/* role that the Amazon ECS container agent and the Docker daemon can assume */
resource "aws_iam_role" "ecs_execution_role" {
  name               = "ecs_task_execution_role"
  assume_role_policy = "${file("${path.module}/policies/ecs-task-execution-role.json")}"
}

resource "aws_iam_role_policy" "ecs_execution_role_policy" {
  name   = "ecs_execution_role_policy"
  policy = "${file("${path.module}/policies/ecs-execution-role-policy.json")}"
  role   = "${aws_iam_role.ecs_execution_role.id}"
}

/*====
ECS service
======*/

/* Security Group for ECS */
resource "aws_security_group" "ecs_service" {
  vpc_id      = "${var.vpc_id}"
  name        = "${var.environment}-ecs-service-sg"
  description = "Allow egress from container"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8
    to_port     = 0
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name        = "${var.environment}-ecs-service-sg"
    Environment = "${var.environment}"
  }
}

// Simply specify the family to find the latest ACTIVE revision in that family
data "aws_ecs_task_definition" "graphql_web" {
  task_definition = "${aws_ecs_task_definition.graphql_web.family}"
}

resource "aws_ecs_service" "graphql_web" {
  name            = "${var.environment}-graphql-web"
  task_definition = "${aws_ecs_task_definition.graphql_web.family}:${max("${aws_ecs_task_definition.graphql_web.revision}", "${data.aws_ecs_task_definition.graphql_web.revision}")}"
  desired_count   = 1
  launch_type     = "FARGATE"
  cluster         = "${aws_ecs_cluster.cluster.id}"
  depends_on      = ["aws_iam_role_policy.ecs_service_role_policy"]

  network_configuration {
    security_groups = ["${var.security_groups_ids}", "${aws_security_group.ecs_service.id}"]
    subnets         = ["${var.subnets_ids}"]
  }

  // Forward traffic from load balancer to the exposed GraphQL container's port
  load_balancer {
    target_group_arn = "${aws_alb_target_group.alb_target_group.arn}"
    container_name   = "graphql"
    container_port   = "${var.graphql_port}"
  }

  depends_on = ["aws_alb_target_group.alb_target_group"]
}

/*====
Auto Scaling for ECS
======*/

resource "aws_iam_role" "ecs_autoscale_role" {
  name               = "${var.environment}_ecs_autoscale_role"
  assume_role_policy = "${file("${path.module}/policies/ecs-autoscale-role.json")}"
}

resource "aws_iam_role_policy" "ecs_autoscale_role_policy" {
  name   = "ecs_autoscale_role_policy"
  policy = "${file("${path.module}/policies/ecs-autoscale-role-policy.json")}"
  role   = "${aws_iam_role.ecs_autoscale_role.id}"
}

resource "aws_appautoscaling_target" "target" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.cluster.name}/${aws_ecs_service.graphql_web.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  role_arn           = "${aws_iam_role.ecs_autoscale_role.arn}"
  min_capacity       = 1
  max_capacity       = 4
}

resource "aws_appautoscaling_policy" "up" {
  name               = "${var.environment}_scale_up"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.cluster.name}/${aws_ecs_service.graphql_web.name}"
  scalable_dimension = "ecs:service:DesiredCount"

  step_scaling_policy_configuration {
    adjustment_type         = "ChangeInCapacity"
    cooldown                = 60
    metric_aggregation_type = "Maximum"

    step_adjustment {
      metric_interval_lower_bound = 0
      scaling_adjustment          = 1
    }
  }

  depends_on = ["aws_appautoscaling_target.target"]
}

resource "aws_appautoscaling_policy" "down" {
  name               = "${var.environment}_scale_down"
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.cluster.name}/${aws_ecs_service.graphql_web.name}"
  scalable_dimension = "ecs:service:DesiredCount"

  step_scaling_policy_configuration {
    adjustment_type         = "ChangeInCapacity"
    cooldown                = 60
    metric_aggregation_type = "Maximum"

    step_adjustment {
      metric_interval_lower_bound = 0
      scaling_adjustment          = -1
    }
  }

  depends_on = ["aws_appautoscaling_target.target"]
}

/* metric used for auto scale */
resource "aws_cloudwatch_metric_alarm" "service_cpu_high" {
  alarm_name          = "${var.environment}_graphql-services_web_cpu_utilization_high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "60"
  statistic           = "Maximum"
  threshold           = "85"

  dimensions {
    ClusterName = "${aws_ecs_cluster.cluster.name}"
    ServiceName = "${aws_ecs_service.graphql_web.name}"
  }

  alarm_actions = ["${aws_appautoscaling_policy.up.arn}"]
  ok_actions    = ["${aws_appautoscaling_policy.down.arn}"]
}
