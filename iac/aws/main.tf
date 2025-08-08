provider "aws" {
  region = var.region
}

resource "aws_ecs_cluster" "main" {
  name = "majesticcoding-cluster"
}

resource "aws_ecs_task_definition" "main" {
  family                   = "majesticcoding-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "512"
  memory                   = "1024"

  container_definitions = jsonencode([{
    name      = "majesticcoding"
    image     = "mattmajestic/majesticcoding:latest"
    essential = true
    portMappings = [{
      containerPort = 8080
      protocol      = "tcp"
    }]
  }])
}

resource "aws_ecs_service" "main" {
  name            = "majesticcoding-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.main.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [var.subnet_id]
    assign_public_ip = true
    security_groups  = [var.security_group_id]
  }
}

variable "region" { default = "us-east-1" }
variable "subnet_id" {}
variable "security_group_id" {}