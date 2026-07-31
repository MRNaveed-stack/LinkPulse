output "db_subnet_group_name" {
  value = aws_db_subnet_group.postgres.name
}

output "rds_security_group_id" {
  value = aws_security_group.rds.id
}

output "db_endpoint" {
  value = split(":", aws_db_instance.postgres.endpoint)[0]
}

output "db_port" {
  value = aws_db_instance.postgres.port
}

output "db_name" {
  value = aws_db_instance.postgres.db_name
}