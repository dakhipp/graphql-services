/*====
Bastion Server
======*/
resource "aws_instance" "bastion_instance" {
  ami           = "${data.aws_ami.ec2-linux.id}"
  instance_type = "t2.micro"

  vpc_security_group_ids = ["${aws_security_group.bastion-sg.id}", "${var.rds_sg}"]
  key_name               = "${aws_key_pair.bastion_key.key_name}"
  subnet_id              = "${var.subnet_id}"

  // User Data is a command that is ran at boot up time, updates, install psql, and stops instance
  user_data = <<-EOF
              #!/bin/bash
              # Install updates and PSQL CLI
              sudo yum upgrade -y ;
              sudo yum install postgresql96-server.x86_64 -y ;
              # Create limited web user
              a="postgresql://${var.psql_root_user}:${var.psql_root_pass}@${var.psql_addr}"
              psql "$a/${var.psql_root_db}" -c "CREATE DATABASE ${var.psql_web_db};"
              psql "$a/${var.psql_web_db}" -c "CREATE USER ${var.psql_web_user} WITH ENCRYPTED PASSWORD '${var.psql_web_pass}';"
              psql "$a/${var.psql_web_db}" -c "GRANT CONNECT ON DATABASE ${var.psql_web_db} TO ${var.psql_web_user};"
              psql "$a/${var.psql_web_db}" -c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE ON TABLES TO ${var.psql_web_user};"
              # Power off machine since it has access to private resources
              sudo poweroff
              EOF

  tags {
    Name        = "bastion_${var.environment}"
    Environment = "${var.environment}"
  }
}

// Gets the most recent version of Amazon linux
data "aws_ami" "ec2-linux" {
  most_recent = true

  filter {
    name   = "name"
    values = ["amzn-ami-*-x86_64-gp2"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "owner-alias"
    values = ["amazon"]
  }
}

// Allows inbound connections on port 22 and all outbound connections
resource "aws_security_group" "bastion-sg" {
  name   = "bastion_security_group"
  vpc_id = "${var.vpc_id}"

  ingress {
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    protocol    = -1
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

// Creates key pair based on var originating from terraform.tfvars
resource "aws_key_pair" "bastion_key" {
  key_name   = "${var.bastion_key_name}"
  public_key = "${var.bastion_public_key}"
}
