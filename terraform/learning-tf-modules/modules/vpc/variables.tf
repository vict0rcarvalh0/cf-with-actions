variable "vpc_cidr_block" {
  description = "CIDR block da VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "enable_dns_support" {
  description = "Habilita suporte a DNS na VPC"
  type        = bool
  default     = true
}

variable "enable_dns_hostnames" {
  description = "Habilita hostnames DNS na VPC"
  type        = bool
  default     = true
}

variable "env" {
    description = "Ambiente"
    type        = string
    default     = "dev"
}

variable "vpc_tags" {
    description = "Tags da VPC"
    type        = map(any)
    default     = { "Company" = "Example"}
}

variable "subnet_cidr_block" {
    description = "CIDR block da subnet"
    type        = list(string)
    default     = ["10.0.0.0/24", "10.0.2.0/24"]
}

variable "azs" {
    description = "Zonas de disponibilidade"
    type        = list(string)
    default     = ["us-east-1a", "us-east-1b"]
}

variable "subnet_tags" {
    description = "Tags da subnet"
    type        = map(any)
    default     = { "Company" = "Example" }
}

 