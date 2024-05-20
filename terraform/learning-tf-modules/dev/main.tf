module "vpc" {
    source = "../modules/vpc"

    vpc_cidr_block = "10.0.0.0/16"
    env = "dev"
    vpc_tags = { "Owner" = "Victor"}
    subnet_cidr_block = ["10.0.0.0/24", "10.0.2.0/24"]
    azs = ["us-east-1a", "us-east-1b"]
    subnet_tags = { "Owner" = "Victor" }
}