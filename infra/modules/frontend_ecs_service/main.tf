resource "aws_ecs_service" "frontend" {
  name            = "${var.project_name}-frontend-service"
  cluster         = var.cluster_id
  task_definition = var.task_definition_arn

  desired_count = 1
  launch_type   = "FARGATE"

  network_configuration {
    subnets = var.subnet_ids

    security_groups = [
      var.ecs_security_group_id
    ]

    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = var.target_group_arn
    container_name   = "frontend"
    container_port   = 80
  }
}