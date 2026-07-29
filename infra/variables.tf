variable "aws_region" {
  description = "AWS Region"
  type        = string
}

variable "project_name" {
  description = "Project Name"
  type        = string
}

variable "db_username" {
  type = string
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "google_client_id" {
  type        = string
  description = "The OAuth 2.0 Client ID for Google Authentication"
}

variable "google_client_secret" {
  type        = string
  description = "The OAuth 2.0 Client Secret for Google Authentication"
  sensitive   = true
}
