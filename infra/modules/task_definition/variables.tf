variable "project_name" {
  type = string
}

variable "execution_role_arn" {
  type = string
}

variable "repository_url" {
  type = string
}

variable "log_group_name" {
  type = string
}

variable "aws_region" {
  type    = string
  default = "ap-south-1"
}

variable "db_host" {
  type = string
}

variable "db_name" {
  type = string
}

variable "db_username" {
  type = string
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "task_role_arn" {
  type = string
}




variable "google_client_id" {
  type        = string
  description = "The OAuth 2.0 Client ID generated from the Google Developer Console"
}

variable "google_client_secret" {
  type        = string
  description = "The OAuth 2.0 Client Secret generated from the Google Developer Console"
  sensitive   = true
}
