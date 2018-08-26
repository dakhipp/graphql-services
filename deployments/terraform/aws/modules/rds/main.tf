/*====
RDS
======*/

// subnets used by RDS
resource "aws_db_subnet_group" "rds_subnet_group" {
  name        = "${var.environment}-rds-subnet-group"
  description = "RDS subnet group"
  subnet_ids  = ["${var.subnet_ids}"]

  tags {
    Environment = "${var.environment}"
  }
}

// Security Group for resources that want to access the Database
resource "aws_security_group" "db_access_sg" {
  vpc_id      = "${var.vpc_id}"
  name        = "${var.environment}-db-access-sg"
  description = "Allow access to RDS"

  tags {
    Name        = "${var.environment}-db-access-sg"
    Environment = "${var.environment}"
  }
}

// Security group for the RDS service, only allows inbound access on port 5432
resource "aws_security_group" "rds_sg" {
  name        = "${var.environment}-rds-sg"
  description = "${var.environment} Security Group"
  vpc_id      = "${var.vpc_id}"

  tags {
    Name        = "${var.environment}-rds-sg"
    Environment = "${var.environment}"
  }

  // Allows traffic from the SG itself
  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    self      = true
  }

  // Allow traffic for TCP 5432
  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = ["${aws_security_group.db_access_sg.id}"]
  }

  // Allow outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_db_instance" "rds" {
  identifier              = "${var.environment}-database"
  allocated_storage       = "${var.allocated_storage}"
  engine                  = "${var.engine}"
  engine_version          = "${var.engine_version}"
  instance_class          = "${var.instance_class}"
  multi_az                = "${var.multi_az}"
  apply_immediately       = "${var.apply_immediately}"
  backup_retention_period = "${var.backup_retention_period}"
  backup_window           = "${var.backup_window}"
  name                    = "${var.psql_db}"
  username                = "${var.psql_user}"
  password                = "${var.psql_pass}"
  db_subnet_group_name    = "${aws_db_subnet_group.rds_subnet_group.id}"
  vpc_security_group_ids  = ["${aws_security_group.rds_sg.id}"]
  skip_final_snapshot     = true

  tags {
    Name        = "${var.environment}-rds"
    Environment = "${var.environment}"
  }
}
