resource "aws_vpc" "main" {
  cidr_block = var.vpc_cidr_block

  enable_dns_support = var.enable_dns_support
  enable_dns_hostnames = var.enable_dns_hostnames

  tags = merge(
    { Name = "${var.env}-main" }, var.vpc_tags
    )
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "internet-gateway-${var.env}"
  }
}

resource "aws_subnet" "public_subnet" {
    count = length(var.subnet_cidr_block)

    vpc_id = aws_vpc.main.id
    cidr_block = var.subnet_cidr_block[count.index]
    availability_zone = var.azs[count.index]
    map_public_ip_on_launch = true
    
    tags = merge(
        { Name = "${var.env}-public-${var.azs[count.index]}" }, var.subnet_tags
    )
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = var.rt_cidr_block
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "public-route-table-${var.env}"
  }
}

resource "aws_route_table_association" "public_subnet_1a" {
  count = length(var.subnet_cidr_block)

  subnet_id      = aws_subnet.public_subnet[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_security_group" "elb" {
  name_prefix = "elb-security-group-${var.env}"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = var.sg_cidr_block
  }

  ingress {
    from_port   = var.ingress_port
    to_port     = var.ingress_port
    protocol    = "tcp"
    cidr_blocks = var.sg_cidr_block
  }

  tags = merge(
        { Name = "${var.env}-elb-sg" }, var.elb_tags
    )
}

resource "aws_security_group" "ec2" {
  name_prefix = "ec2-security-group-${var.env}"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    security_groups = [aws_security_group.elb.id]
  }

  ingress {
    from_port   = var.ingress_port
    to_port     = var.ingress_port
    protocol    = "tcp"
    cidr_blocks = var.sg_cidr_block
  }

  tags = merge(
        { Name = "${var.env}-ec2-sg" }, var.ec2_tags
    )
}

resource "aws_instance" "backend_instance" {
  count                      = length(var.subnet_cidr_block)

  ami                        = data.aws_ssm_parameter.ami_id.value
  instance_type              = var.instance_type
  subnet_id                  = aws_subnet.public_subnet[count.index].id
  key_name                   = var.key_pair_name
  vpc_security_group_ids     = [aws_security_group.ec2.id]
  user_data = base64encode(<<-EOF
                #!/bin/bash

                yum update -y
                amazon-linux-extras install docker -y
                systemctl start docker

                touch home/ec2-user/script.sh &>> $LOG_FILE
              EOF
            )

  tags = {
    Name = "backend-1b-${var.env}"
  }
}

resource "tls_private_key" "pk" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "kp" {
  key_name   = var.key_pair_name
  public_key = tls_private_key.pk.public_key_openssh

  provisioner "local-exec" {
    command = "echo '${tls_private_key.pk.private_key_pem}' > ./'${var.key_pair_name}'.pem"
  }
}

resource "aws_lb_target_group" "ec2" {
  name     = var.lb_tg_name
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id

  health_check {
    interval            = 30
    protocol            = "HTTP"
    timeout             = 15
    healthy_threshold   = 5
    unhealthy_threshold = 3
    matcher             = "200"
  }
}

resource "aws_lb_target_group_attachment" "ec2_instance" {
  count            = length(var.subnet_cidr_block)

  target_group_arn = aws_lb_target_group.ec2.arn
  target_id        = aws_instance.backend_instance[count.index].id
  port             = 80
}

resource "aws_lb" "elb" {
  name               = "elb-${var.env}"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.elb.id]
  subnets            = aws_subnet.public_subnet[*].id

  enable_deletion_protection = false

  tags = {
    Name = "elb-${var.env}"
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.elb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.ec2.arn
  }
}