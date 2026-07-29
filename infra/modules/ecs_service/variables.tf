# modules/ecs_service/variables.tf

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
}

# ⭐ ADD THESE NEW VARIABLES:
variable "cluster_id" {
  description = "ID of the ECS cluster"
  type        = string
}

variable "task_definition_arn" {
  description = "ARN of the ECS task definition"
  type        = string
}

variable "subnet_ids" {
  description = "List of subnet IDs for the ECS service"
  type        = list(string)
}

variable "ecs_security_group_id" {
  description = "Security group ID for ECS tasks"
  type        = string
}

variable "target_group_arn" {
  description = "ARN of the target group for load balancer"
  type        = string
}

# Optional: If you want to be more flexible
variable "desired_count" {
  description = "Number of ECS service instances"
  type        = number
  default     = 1
}

variable "container_name" {
  description = "Container name in the task definition"
  type        = string
  default     = "backend"
}

variable "container_port" {
  description = "Container port for the application"
  type        = number
  default     = 8080
}

variable "assign_public_ip" {
  description = "Assign public IP to ECS tasks"
  type        = bool
  default     = true
}