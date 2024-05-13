output "key_pair_name" {
    value = data.aws_key_pair.my_key_web.key_name
    sensitive = true
    description = "Value of the key pair name"
}

output "sg_id" {
    value = data.aws_security_group.sg_web.id
    sensitive = false
    description = "Value of the security group id"
}