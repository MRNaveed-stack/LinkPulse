variable "project_name" {
  type = string
}

variable "task_role_arn" {
  type        = string
  description = "The ARN of the IAM role that allows the ECS tasks to make calls to AWS API services"
}
