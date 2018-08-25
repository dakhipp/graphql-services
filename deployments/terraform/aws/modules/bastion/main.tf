/*====
Bastion Server
======*/
resource "aws_instance" "bastion_instance" {
  ami           = "${data.aws_ami.ec2-linux.id}"
  instance_type = "t2.micro"

  vpc_security_group_ids = ["${aws_security_group.bastion-sg.id}", "${var.rds_sg}"]
  key_name               = "${aws_key_pair.bastion_key.key_name}"
  subnet_id              = "${var.subnet_id}"

  tags {
    Name        = "bastion-${var.environment}"
    Environment = "${var.environment}"
  }
}

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

resource "aws_security_group" "bastion-sg" {
  name   = "bastion-security-group"
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

resource "aws_key_pair" "bastion_key" {
  key_name   = "${var.bastion_key_name}"
  public_key = "${var.bastion_public_key}"
}
