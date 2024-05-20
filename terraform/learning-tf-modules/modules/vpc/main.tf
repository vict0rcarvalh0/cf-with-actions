// Parâmetros dos resources viram variáveis
resource "aws_vpc" "main" {
  cidr_block = var.vpc_cidr_block

  enable_dns_support = var.enable_dns_support
  enable_dns_hostnames = var.enable_dns_hostnames

  tags = merge(
    { Name = "${var.env}-main" }, var.vpc_tags
    )
}

resource "aws_subnet" "private_us_east_1a" {
    count = length(var.subnet_cidr_block)

    vpc_id = aws_vpc.main.id
    cidr_block = var.subnet_cidr_block[count.index]
    availability_zone = var.azs[count.index]
    
    tags = merge(
        { Name = "${var.env}-private-${var.azs[count.index]}" }
    )
}