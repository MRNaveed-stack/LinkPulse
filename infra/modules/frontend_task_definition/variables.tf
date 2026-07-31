variable "project_name" {
  type = string
}

variable "repository_url" {
  type = string
}

variable "execution_role_arn" {
  type = string
}

variable "log_group_name" {
  type = string
}

variable "aws_region" {
  type    = string
  default = "ap-south-1"
}