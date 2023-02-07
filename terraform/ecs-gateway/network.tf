resource "aws_lb" "ecs-nlb" {
  name               = "ecs-nlb-${random_id.id.hex}"
  load_balancer_type = "network"
  ip_address_type    = "ipv4"

  subnet_mapping {
    subnet_id = data.aws_subnet.public-subnet.id
  }
}

resource "aws_lb_listener" "ecs-nlb-listener" {
  load_balancer_arn = aws_lb.ecs-nlb.arn
  port              = "5000"
  protocol          = "TCP_UDP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.ecs-target-group.arn
  }
}

resource "aws_lb_target_group" "ecs-target-group" {
  name        = "ecs-target-group-${random_id.id.hex}"
  port        = 5000
  protocol    = "TCP_UDP"
  target_type = "ip"
  vpc_id      = data.aws_vpc.vpc.id
}


data "aws_vpc" "vpc" {
  id = var.vpc_id
}

data "aws_subnet" "public-subnet" {
  id = var.public_subnet_id
}

data "aws_subnet" "private-subnet" {
  id = var.private_subnet_id
}

resource "aws_security_group" "sdm_sg" {
  name        = "sdm-sg-${random_id.id.hex}"
  vpc_id      = data.aws_vpc.vpc.id
  description = "Sec Group for sdm"

}

resource "aws_security_group_rule" "sdm_ingress_rules" {
  count = length(var.sg_sdm_ingress_rules)
  type              = "ingress"
  from_port         = var.sg_sdm_ingress_rules[count.index].from_port
  to_port           = var.sg_sdm_ingress_rules[count.index].to_port
  protocol          = var.sg_sdm_ingress_rules[count.index].protocol
  cidr_blocks       = [var.sg_sdm_ingress_rules[count.index].cidr_block]
  description       = var.sg_sdm_ingress_rules[count.index].description
  security_group_id = aws_security_group.sdm_sg.id

}

resource "aws_security_group_rule" "egress_sdm_rules" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.sdm_sg.id

}
