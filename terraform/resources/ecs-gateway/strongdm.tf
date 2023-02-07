resource "sdm_node" "sdm-ecs-gateway-01" {
  gateway {
    name           = "sdm-ecs-${random_id.id.hex}"
    listen_address = "${aws_lb.ecs-nlb.dns_name}:5000"
  }
}