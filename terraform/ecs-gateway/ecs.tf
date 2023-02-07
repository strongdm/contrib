resource "aws_ecs_task_definition" "sdm-ecs-task-definition" {
  family                   = "sdm-ecs${random_id.id.hex}-terraform"
  task_role_arn            = aws_iam_role.sdm_ecs_task_role.arn
  execution_role_arn       = aws_iam_role.sdm_ecs_task_execution_role.arn
  network_mode             = "awsvpc"
  cpu                      = "512"
  memory                   = "1024"
  requires_compatibilities = ["FARGATE"]
  container_definitions = jsonencode(
    [
      {
        "image" : "quay.io/sdmrepo/relay",
        "name" : "sdm-container",
        "memory" : 1024,
        "portmappings" : [{
          "containerport" : 5000
          }
        ],
        "linuxParameters" : {
          "initProcessEnabled" : true
        },
        "logConfiguration": {
            "logDriver": "awslogs",
            "options": {
                "awslogs-group": "${aws_cloudwatch_log_group.node.name}",
                "awslogs-region": "${var.region}", 
                "awslogs-stream-prefix": "ecsgateway"
            }
        },
        "environment" : [
          {
            "name" : "SDM_RELAY_TOKEN",
            "value" : "${sdm_node.sdm-ecs-gateway-01.gateway[0].token}"
          },
          {
            "name" : "SDM_DOCKERIZED",
            "value" : "true"
          }
        ]
      }

  ])
  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }
}

resource "aws_ecs_service" "sdm-gateway" {
  name                               = "sdm-gateway-${random_id.id.hex}"
  cluster                            = aws_ecs_cluster.sdm-ecs-cluster.id
  task_definition                    = aws_ecs_task_definition.sdm-ecs-task-definition.arn
  desired_count                      = 1
  depends_on                         = [aws_iam_role.sdm_ecs_task_role]
  deployment_maximum_percent         = 100
  deployment_minimum_healthy_percent = 1
  enable_execute_command             = true

  load_balancer {
    target_group_arn = aws_lb_target_group.ecs-target-group.arn
    container_name   = "sdm-container"
    container_port   = 5000
  }
  network_configuration {
    subnets         = [data.aws_subnet.private-subnet.id]
    security_groups = [aws_security_group.sdm_sg.id]
  }
  capacity_provider_strategy {
    base              = 1
    weight            = 1
    capacity_provider = "FARGATE"
  }
}

resource "aws_ecs_cluster" "sdm-ecs-cluster" {
  name = "sdm-gateway-cluster"
}

# Create Cloudwatch log group for the container
resource "aws_cloudwatch_log_group" "node" {
  name = "ecs-${random_id.id.hex}-gateway"
}
