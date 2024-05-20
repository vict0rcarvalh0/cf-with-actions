locals {
    vpc_cidr_prefix = "10.1"
}

module "vpc" {
    source = "../modules/vpc"

    vpc_cidr_block = "${local.vpc_cidr_prefix}.0.0/16"
    env = "prod"
    vpc_tags = { "Owner" = "Victor"}
    subnet_cidr_block = ["${local.vpc_cidr_prefix}.0.0/24", "${local.vpc_cidr_prefix}.1.0/24"]
    azs = ["us-east-1a", "us-east-1b"]
    subnet_tags = { "Owner" = "Victor" }
}