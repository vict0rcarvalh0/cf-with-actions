output "vpc_id" {
  value = aws_vpc.main.id
}

output "subnet_id_1a" {
  value = aws_subnet.public_subnet_1a.id
}

output "subnet_id_1b" {
  value = aws_subnet.public_subnet_1b.id
}

output "backend_instance_1_id" {
  value = aws_instance.backend_instance_1.id
}

output "backend_instance_2_id" {
  value = aws_instance.backend_instance_2.id
}

output "db_instance_address" {
  value = aws_db_instance.default.address
}