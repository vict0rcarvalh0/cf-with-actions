output "vpc_id" {
    description = "ID da VPC"
    value       = aws_vpc.this.id
}

output "subnet_ids" {
    description = "IDs das subnets"
    value       = aws_subnet.this[*].id
}