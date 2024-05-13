data "aws_key_pair" "my_key_web" {
    key_name         = "my-key-pair"
  include_public_key = true
}

data "aws_security_group" "sg_web" {
  name = "sg-web"
}

data "aws_subnet" "subnet_web" {
    filter {
        name = "tag:Name"
        values = ["subnet-web-a"]
    }
}

locals {
    tags = {
        Name = var.bucket_name
        Owner = "DevOps Team"
    }
}

resource "aws_s3_bucket" "bucket_s3" {
    bucket = var.bucket_name
    tags   = local.tags
}