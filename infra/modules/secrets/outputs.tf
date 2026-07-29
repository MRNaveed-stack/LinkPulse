output "secret_arn" {
  value = aws_secretsmanager_secret.database.arn
}

output "secret_name" {
  value = aws_secretsmanager_secret.database.name
}