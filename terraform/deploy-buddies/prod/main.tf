locals {
    cidr_prefix = "10.0"
    cidr_block = "0.0.0.0/0"
}

module "resources" {
    source = "../modules/resources"

    vpc_cidr_block = "${local.cidr_prefix}.0.0/16"
    env = "prod"
    vpc_tags = { "Company" = "DeployBuddy"}
    subnet_cidr_block = ["${local.cidr_prefix}.1.0/24", "${local.cidr_prefix}.2.0/24"]
    azs = ["us-east-1a", "us-east-1b"]
    subnet_tags = { "Company" = "DeployBuddy"}
    sg_cidr_block = ["${local.cidr_block}"]
    ingress_port = 8080
    elb_tags = { "Company" = "DeployBuddy"}
    ec2_tags = { "Company" = "DeployBuddy"}
    instance_type = "t2.micro"
    rt_cidr_block = "${local.cidr_block}"
    key_pair_name = "DEPLOYBUDDY-KP"
}