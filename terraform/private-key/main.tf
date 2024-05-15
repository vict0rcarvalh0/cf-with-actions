provider "aws" {
  region = "us-east-1"
}

variable "key_pair_name" {
  description = "key_pair_name"
  type        = string
}

resource "tls_private_key" "pk" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "kp" {
  key_name   = var.key_pair_name
  public_key = tls_private_key.pk.public_key_openssh

  provisioner "local-exec" {
    command = "echo '${tls_private_key.pk.private_key_pem}' > ./KP-TESTE.pem"
  }
}