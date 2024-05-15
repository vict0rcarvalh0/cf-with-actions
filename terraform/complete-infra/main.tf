resource "aws_vpc" "main" {
  cidr_block       = "10.0.0.0/16"
  enable_dns_support = true
  enable_dns_hostnames = true

  tags = {
    Name = "vpc-${var.environment_type}"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "internet-gateway-${var.environment_type}"
  }
}

resource "aws_subnet" "public_subnet_1a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true

  tags = {
    Name = "public-subnet-1a-${var.environment_type}"
  }
}

resource "aws_subnet" "public_subnet_1b" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.3.0/24"
  availability_zone       = "us-east-1b"
  map_public_ip_on_launch = true

  tags = {
    Name = "public-subnet-ab-${var.environment_type}"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "public-route-table-${var.environment_type}"
  }
}

resource "aws_route_table_association" "public_subnet_1a" {
  subnet_id      = aws_subnet.public_subnet_1a.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "public_2" {
  subnet_id      = aws_subnet.public_subnet_1b.id
  route_table_id = aws_route_table.public.id
}

resource "aws_security_group" "elb" {
  name_prefix = "elb-security-group-${var.environment_type}"
  description = "ELB Security Group"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "elb-security-group-${var.environment_type}"
  }
}

resource "aws_security_group" "ec2" {
  name_prefix = "ec2-security-group-${var.environment_type}"
  description = "EC2 Security Group"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    security_groups = [aws_security_group.elb.id]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "ec2-security-group-${var.environment_type}"
  }
}

resource "aws_instance" "backend_instance_1" {
  ami                         = data.aws_ssm_parameter.ami_id.value
  instance_type              = "t2.micro"
  subnet_id                  = aws_subnet.public_subnet_1a.id
  key_name                   = var.key_pair_name
  vpc_security_group_ids     = [aws_security_group.ec2.id]
#   user_data                  = file("userdata.sh")

  tags = {
    Name = "backend-1-${var.environment_type}"
  }
}

resource "aws_instance" "backend_instance_2" {
  ami                         = data.aws_ssm_parameter.ami_id.value
  instance_type              = "t2.micro"
  subnet_id                  = aws_subnet.public_subnet_1b.id
  key_name                   = var.key_pair_name
  vpc_security_group_ids     = [aws_security_group.ec2.id]
#   user_data                  = file("userdata.sh")

  tags = {
    Name = "backend-2-${var.environment_type}"
  }
}

resource "aws_lb_target_group" "ec2" {
  name     = "EC2TargetGroup"
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

resource "aws_lb_target_group_attachment" "ec2_instance_1" {
  target_group_arn = aws_lb_target_group.ec2.arn
  target_id        = aws_instance.backend_instance_1.id
  port             = 80
}

resource "aws_lb_target_group_attachment" "ec2_instance_2" {
  target_group_arn = aws_lb_target_group.ec2.arn
  target_id        = aws_instance.backend_instance_2.id
  port             = 80
}

resource "aws_lb" "elb" {
  name               = "elb-${var.environment_type}"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.elb.id]
  subnets            = [aws_subnet.public_subnet_1a.id, aws_subnet.public_subnet_1b.id]

  enable_deletion_protection = false

  tags = {
    Name = "elb-${var.environment_type}"
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

resource "aws_db_subnet_group" "main" {
  name       = "main-subnet-group"
  subnet_ids = [aws_subnet.public_subnet_1a.id, aws_subnet.public_subnet_1b.id]

  tags = {
    Name = "main-subnet-group"
  }
}

resource "aws_db_instance" "default" {
  allocated_storage      = 5
  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.ec2.id]
  identifier             = var.db_instance_identifier
  instance_class         = "db.t3.micro"
  engine                 = "postgres"
  username               = var.db_username
  password               = var.db_password
  skip_final_snapshot    = true

  tags = {
    Name = "rds-${var.environment_type}"
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

resource "aws_launch_configuration" "app_launch_configuration" {
  name_prefix   = "app-launch-configuration-"
  image_id      = var.ami_ssm_parameter
  instance_type = "t2.micro"
  key_name      = aws_key_pair.deployer_key_pair.key_name

  security_groups = [aws_security_group.ec2.id]

  user_data = <<-EOF
                #!/bin/bash
                echo "Hello, Terraform!" > /var/www/html/index.html
                EOF
}

resource "aws_autoscaling_group" "app_autoscaling_group" {
  launch_configuration = aws_launch_configuration.app_launch_configuration.id
  min_size             = 1
  max_size             = 3
  desired_capacity     = 2

  vpc_zone_identifier = [aws_subnet.public_subnet_1a.id, aws_subnet.public_subnet_1b.id]

  tag {
    key                 = "Name"
    value               = "asg-instance"
    propagate_at_launch = true
  }

  wait_for_capacity_timeout = "15m"

  lifecycle {
    ignore_changes = [id]
  }
}
