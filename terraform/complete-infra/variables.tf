variable "environment_type" {
  description = "Specify the Environment type of the stack."
  type        = string
  validation {
    condition     = contains(["dev", "uat", "prod"], var.environment_type)
    error_message = "Environment type must be either 'dev', 'uat', or 'prod'."
  }
}

variable "key_pair_name" {
  description = "The name of an existing Amazon EC2 key pair in this region to use to SSH into the Amazon EC2 instances."
  type        = string
}

variable "db_instance_identifier" {
  description = "The identifier for the RDS instance."
  type        = string
  default     = "postgres-db"
}

variable "db_username" {
  description = "The username for the RDS instance."
  type        = string
  default     = "postgres"
}

variable "db_password" {
  description = "The password for the RDS instance."
  type        = string
  sensitive   = true
}

variable "ami_ssm_parameter" {
  description = "The SSM parameter for the AMI ID."
  type        = string
  default     = "/aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2"
}
