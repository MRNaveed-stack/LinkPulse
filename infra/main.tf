resource "aws_s3_bucket" "linkpulse_uploads" {
  bucket = "${var.project_name}-mnq-uploads"
  tags = {
    Project = "LinkPulse"
    Managed = "Terraform"
  }
}

module "uploads_bucket" {
  source = "./modules/s3"

  project_name = var.project_name
}

module "ecr" {
  source = "./modules/ecr"

  project_name = var.project_name
}

module "vpc" {
  source = "./modules/vpc"

  project_name = var.project_name
}

module "security_groups" {
  source = "./modules/security_groups"

  project_name = var.project_name
  vpc_id       = module.vpc.vpc_id
}


module "ecs" {
  source = "./modules/ecs"

  project_name  = var.project_name
  task_role_arn = module.ecs_task_role.task_role_arn


}


module "cloudwatch" {
  source = "./modules/cloudwatch"

  project_name = var.project_name
}

module "ecs_task_execution_role" {
  source = "./modules/ecs_task_execution_role"

  project_name = var.project_name
}


module "task_definition" {
  source = "./modules/task_definition"

  project_name  = var.project_name
  task_role_arn = module.ecs_task_role.task_role_arn

  execution_role_arn   = module.ecs_task_execution_role.execution_role_arn
  google_client_id     = var.google_client_id
  google_client_secret = var.google_client_secret

  repository_url = module.ecr.repository_url

  log_group_name = module.cloudwatch.log_group_name

  db_host = module.rds.db_endpoint

  db_name = module.rds.db_name

  db_username = var.db_username

  db_password = var.db_password

  jwt_secret           = var.jwt_secret
  google_redirect_url  = "http://${module.alb.alb_dns_name}/api/auth/google/callback"
  frontend_url         = "http://${module.alb.alb_dns_name}"
}


module "alb" {

  source = "./modules/alb"

  project_name = var.project_name

  vpc_id = module.vpc.vpc_id

  public_subnet_ids = module.vpc.public_subnet_ids

  alb_security_group_id = module.security_groups.alb_security_group_id

}


# root main.tf - CORRECTED VERSION

module "ecs_service" {
  source = "./modules/ecs_service"

  project_name          = var.project_name
  cluster_id            = module.ecs.cluster_id
  task_definition_arn   = module.task_definition.task_definition_arn
  subnet_ids            = module.vpc.public_subnet_ids
  ecs_security_group_id = module.security_groups.ecs_security_group_id
  target_group_arn      = module.alb.target_group_arn

  depends_on = [
    module.alb
  ]
}

module "rds" {

  source = "./modules/rds"

  project_name = var.project_name

  private_subnet_ids = module.vpc.private_subnet_ids

  vpc_id = module.vpc.vpc_id

  ecs_security_group_id = module.security_groups.ecs_security_group_id

  db_username = var.db_username

  db_password = var.db_password

}


module "secrets" {

  source = "./modules/secrets"

  project_name = var.project_name

  db_host = module.rds.db_endpoint

  db_name = module.rds.db_name

  db_username = var.db_username

  db_password = var.db_password

}

module "ecs_task_role" {

  source = "./modules/ecs_task_role"

  project_name = var.project_name

  secret_arn = module.secrets.secret_arn
}

module "frontend_task_definition" {
  source = "./modules/frontend_task_definition"

  project_name       = var.project_name
  repository_url     = module.ecr.frontend_repository_url
  execution_role_arn = module.ecs_task_execution_role.execution_role_arn
  log_group_name     = module.cloudwatch.log_group_name
}


module "frontend_ecs_service" {
  source = "./modules/frontend_ecs_service"

  project_name = var.project_name

  cluster_id = module.ecs.cluster_id

  task_definition_arn = module.frontend_task_definition.task_definition_arn

  subnet_ids = module.vpc.public_subnet_ids

  ecs_security_group_id = module.security_groups.ecs_security_group_id

  target_group_arn = module.alb.frontend_target_group_arn

  depends_on = [
    module.alb
  ]
}