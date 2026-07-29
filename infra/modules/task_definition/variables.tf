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
