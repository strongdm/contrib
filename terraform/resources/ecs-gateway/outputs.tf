output "gateway-token" {
  value     = sdm_node.sdm-ecs-gateway-01.gateway[0].token
  sensitive = true
}

output "nlb_address" {
  value = aws_lb.ecs-nlb.dns_name
  }