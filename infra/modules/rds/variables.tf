variable "project_name" {
  type = string
}

variable "private_subnet_ids" {
  type = list(string)
}

variable "vpc_id" {
  type = string
}

variable "ecs_security_group_id" {
  type = string
}

variable "db_name" {
  type    = string
  default = "linkpulse"
}

variable "db_username" {
  type = string
}

variable "db_password" {
  type      = string
  sensitive = true
}