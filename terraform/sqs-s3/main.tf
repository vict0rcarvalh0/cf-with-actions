provider "aws" {
  region  = "${var.region}"
}

resource "aws_s3_bucket" "s3_bucket" {
  bucket = "bucket-with-sqs"
    
  tags = {
    Environment = "${var.environment}"
  }
}

resource "aws_sqs_queue" "sqs-blog" {
  name = "sqs-posts"
  delay_seconds = 90
  tags = {
    Environment = "${var.environment}"
  }
}