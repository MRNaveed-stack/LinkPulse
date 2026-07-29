resource "aws_secretsmanager_secret" "database" {

  name = "${var.project_name}-database"

  tags = {
    Project = var.project_name
  }

}

resource "aws_secretsmanager_secret_version" "database" {

  secret_id = aws_secretsmanager_secret.database.id

  secret_string = jsonencode({

    host     = var.db_host

    port     = 5432

    database = var.db_name

    username = var.db_username

    password = var.db_password

  })

}

