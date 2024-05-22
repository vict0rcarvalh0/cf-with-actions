variable "vpc_cidr_block" {
  description = "CIDR block da VPC"
  type        = string
  default     = "10.1.0.0/16"
}

variable "enable_dns_support" {
  description = "Habilita suporte a DNS na VPC"
  type        = bool
  default     = true
}

variable "enable_dns_hostnames" {
  description = "Habilita hostnames DNS na VPC"
  type        = bool
  default     = true
}

variable "env" {
    description = "Ambiente"
    type        = string
    default     = "dev"
}

variable "vpc_tags" {
    description = "Tags da VPC"
    type        = map(any)
    default     = { "Company" = "DeployBuddy"}
}

variable "subnet_cidr_block" {
    description = "CIDR block da subnet"
    type        = list(string)
    default     = ["10.1.0.0/24", "10.1.2.0/24"]
}

variable "azs" {
    description = "Zonas de disponibilidade"
    type        = list(string)
    default     = ["us-east-1a", "us-east-1b"]
}

variable "subnet_tags" {
    description = "Tags da subnet"
    type        = map(any)
    default     = { "Company" = "DeployBuddy" }
}

variable "elb_tags" {
    description = "Tags do ELB"
    type        = map(any)
    default     = { "Company" = "DeployBuddy" }
}

variable "ec2_tags" {
    description = "Tags da EC2"
    type        = map(any)
    default     = { "Company" = "DeployBuddy" }
}

variable "instance_type" {
    description = "Tipo de instância"
    type        = string
    default     = "t2.micro"
}

variable "key_pair_name" {
  description = "The name of an existing Amazon EC2 key pair in this region to use to SSH into the Amazon EC2 instances."
  type        = string
}

variable "lb_tg_name" {
    description = "Nome do target group"
    type        = string
    default     = "deploybuddy-load-balancer-tg"
}

variable "ami_ssm_parameter" {
  description = "The SSM parameter for the AMI ID."
  type        = string
  default     = "/aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2"
}

variable "sg_cidr_block" {
    description = "CIDR block da subnet"
    type        = list(string)
    default     = ["0.0.0.0/0"]
}

variable "rt_cidr_block" {
    description = "CIDR block da subnet"
    type        = string
    default     = "10.1.0.0/16"
}

variable "ingress_port" {
    description = "Porta de entrada"
    type        = number
    default     = 8080
} 