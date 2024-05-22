output "vpc_id" {
    description = "ID da VPC"
    value       = aws_vpc.main.id
}

output "subnet_ids" {
    description = "IDs das subnets"
    value       = [for subnet in aws_subnet.public_subnet : subnet.id]
}

output "backend_instance_ids" {
    description = "IDs das instâncias EC2"
    value       = [for instance in aws_instance.backend_instance : instance.id]
}
