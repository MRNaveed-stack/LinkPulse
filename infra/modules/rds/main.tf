resource "aws_db_subnet_group" "postgres" {

  name = "${var.project_name}-db-subnet-group"

  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.project_name}-db-subnet-group"
  }

}

resource "aws_security_group" "rds" {

  name = "${var.project_name}-rds-sg"

  description = "RDS Security Group"

  vpc_id = var.vpc_id

  ingress {

    from_port = 5432
    to_port   = 5432

    protocol = "tcp"

    security_groups = [
      var.ecs_security_group_id
    ]
  }

  egress {

    from_port = 0
    to_port   = 0

    protocol = "-1"

    cidr_blocks = [
      "0.0.0.0/0"
    ]
  }
}

resource "aws_db_instance" "postgres" {

  identifier = "${var.project_name}-postgres"

  engine = "postgres"

  engine_version = "17"

  instance_class = "db.t3.micro"

  allocated_storage = 20

  storage_type = "gp3"

  db_name = var.db_name

  username = var.db_username

  password = var.db_password

  db_subnet_group_name = aws_db_subnet_group.postgres.name

  vpc_security_group_ids = [
    aws_security_group.rds.id
  ]

  publicly_accessible = false

  multi_az = false

  backup_retention_period = 1

  deletion_protection = false

  skip_final_snapshot = true

  storage_encrypted = true

  tags = {
    Project = var.project_name
  }

}