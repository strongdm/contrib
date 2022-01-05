# Locals

locals {
  required_tags = {
    ExpiryDate = "2021-12-16"
  }
}

# VPC

resource "aws_vpc" "Main" {
  cidr_block       = var.main_vpc_cidr
  instance_tenancy = "default"
  tags             = merge(local.required_tags, var.resource_tags)
}

resource "aws_internet_gateway" "IGW" {
  vpc_id = aws_vpc.Main.id
}

resource "aws_subnet" "publicsubnets" {
  vpc_id     = aws_vpc.Main.id
  cidr_block = var.public_subnets
  tags       = merge(local.required_tags, var.resource_tags)
}

resource "aws_subnet" "privatesubnets" {
  vpc_id     = aws_vpc.Main.id
  cidr_block = var.private_subnets
  tags       = merge(local.required_tags, var.resource_tags)
}

resource "aws_route_table" "PublicRT" {
  vpc_id = aws_vpc.Main.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.IGW.id
  }
  tags = merge(local.required_tags, var.resource_tags)
}

resource "aws_route_table" "PrivateRT" {
  vpc_id = aws_vpc.Main.id
  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.NATgw.id
  }
  tags = merge(local.required_tags, var.resource_tags)
}

resource "aws_route_table_association" "PublicRTassociation" {
  subnet_id      = aws_subnet.publicsubnets.id
  route_table_id = aws_route_table.PublicRT.id
}

resource "aws_route_table_association" "PrivateRTassociation" {
  subnet_id      = aws_subnet.privatesubnets.id
  route_table_id = aws_route_table.PrivateRT.id
}

resource "aws_eip" "natIP" {
  vpc = true
}

resource "aws_nat_gateway" "NATgw" {
  allocation_id = aws_eip.natIP.id
  subnet_id     = aws_subnet.publicsubnets.id
}

# AMI

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

# Security Groups

resource "aws_security_group" "gateway_sg" {

  name        = "sdm_gateway_sg"
  vpc_id      = aws_vpc.Main.id
  description = "Sec Group for SDM Gateway"

}

resource "aws_security_group_rule" "ingress_rules" {

  count = length(var.gateway_sg_ingress_rules)

  type              = "ingress"
  from_port         = var.gateway_sg_ingress_rules[count.index].from_port
  to_port           = var.gateway_sg_ingress_rules[count.index].to_port
  protocol          = var.gateway_sg_ingress_rules[count.index].protocol
  cidr_blocks       = [var.gateway_sg_ingress_rules[count.index].cidr_block]
  description       = var.gateway_sg_ingress_rules[count.index].description
  security_group_id = aws_security_group.gateway_sg.id

}

resource "aws_security_group_rule" "egress_rules" {

  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.gateway_sg.id

}

resource "aws_security_group" "relay_sg" {

  name        = "sdm_relay_sg"
  vpc_id      = aws_vpc.Main.id
  description = "Sec Group for SDM Gateway"

}

resource "aws_security_group_rule" "relay_ingress_rule" {

  type                     = "ingress"
  from_port                = 22
  to_port                  = 22
  protocol                 = "tcp"
  description              = "ssh"
  source_security_group_id = aws_security_group.gateway_sg.id
  security_group_id        = aws_security_group.relay_sg.id

}

resource "aws_security_group_rule" "relay_egress_rules" {

  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.relay_sg.id

}

resource "aws_security_group" "psql_sg" {

  name        = "sdm_psql_sg"
  vpc_id      = aws_vpc.Main.id
  description = "Sec Group for PSQL"

}

resource "aws_security_group_rule" "psql_ingress_rules" {

  count = length(var.psql_sg_ingress_rules)

  type              = "ingress"
  from_port         = var.psql_sg_ingress_rules[count.index].from_port
  to_port           = var.psql_sg_ingress_rules[count.index].to_port
  protocol          = var.psql_sg_ingress_rules[count.index].protocol
  cidr_blocks       = [var.psql_sg_ingress_rules[count.index].cidr_block]
  description       = var.psql_sg_ingress_rules[count.index].description
  security_group_id = aws_security_group.psql_sg.id

}

resource "aws_security_group_rule" "psql_egress_rules" {

  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.psql_sg.id

}

# SDM Gateway Instance

resource "sdm_node" "gateway" {
  gateway {
    name           = "sdm-gateway-01"
    listen_address = "${aws_eip.gateway.public_ip}:5000"
    bind_address   = "0.0.0.0:5000"
  }
}

output "gateway_token" {
  value     = sdm_node.gateway.gateway[0].token
  sensitive = true
}

resource "aws_eip" "gateway" {
  network_interface = aws_network_interface.gateway.id
}

resource "aws_network_interface" "gateway" {
  subnet_id       = aws_subnet.publicsubnets.id
  security_groups = [aws_security_group.gateway_sg.id]

  tags = merge(local.required_tags, var.resource_tags)
}

resource "aws_instance" "gateway" {

  ami           = data.aws_ami.ubuntu.image_id
  instance_type = "t2.micro"
  key_name      = aws_key_pair.terraform_key.key_name

  user_data = templatefile("${path.module}/template/sdm_gateway_install/install.sh.tpl", { SDM_GATEWAY_TOKEN = "${sdm_node.gateway.gateway[0].token}" })
  tags      = merge(local.required_tags, var.resource_tags)

  network_interface {
    network_interface_id = aws_network_interface.gateway.id
    device_index         = 0
  }

}

# SDM Relay Instance

resource "sdm_node" "relay" {
  relay {
    name = "sdm-relay-01"
  }
}

output "relay_token" {
  value     = sdm_node.relay.relay[0].token
  sensitive = true
}

resource "aws_instance" "relay" {
  depends_on             = [aws_nat_gateway.NATgw]
  ami                    = data.aws_ami.ubuntu.image_id
  instance_type          = "t2.micro"
  subnet_id              = aws_subnet.privatesubnets.id
  vpc_security_group_ids = [aws_security_group.relay_sg.id]
  key_name               = aws_key_pair.terraform_key.key_name

  user_data = templatefile("${path.module}/template/sdm_relay_install/install.sh.tpl", { SDM_RELAY_TOKEN = "${sdm_node.relay.relay[0].token}", SSH_PUB_KEY = "${data.sdm_ssh_ca_pubkey.ssh_pubkey_query.public_key}" })
  tags      = merge(local.required_tags, var.resource_tags)

}

resource "sdm_resource" "relay_ssh" {
  ssh_cert {
    name     = "sdm-relay-ssh"
    username = "ubuntu"
    hostname = aws_instance.relay.private_ip
    port     = 22
    tags     = merge(local.required_tags, var.resource_tags)
  }
}

# PSQL Instance

resource "aws_instance" "psql" {
  depends_on             = [aws_nat_gateway.NATgw]
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = "t3.small"
  vpc_security_group_ids = [aws_security_group.psql_sg.id]
  subnet_id              = aws_subnet.privatesubnets.id
  user_data              = templatefile("${path.module}/template/psql/install.sh.tpl", { SSH_PUB_KEY = "${data.sdm_ssh_ca_pubkey.ssh_pubkey_query.public_key}" })
  tags                   = merge(local.required_tags, var.resource_tags)
}

resource "sdm_resource" "psql_admin" {
  postgres {
    name     = "sdm-psql-admin"
    hostname = aws_instance.psql.private_ip
    database = "dvdrental"
    username = "postgres"
    password = "notastrongpassword123"
    port     = 5432

    tags = merge(local.required_tags, var.resource_tags)
  }
}

resource "sdm_role_grant" "admin_grant_psql_admin" {
  role_id     = sdm_role.admins.id
  resource_id = sdm_resource.psql_admin.id
}

resource "sdm_resource" "psql_ssh" {
  ssh_cert {
    name     = "sdm-psql-ssh"
    username = "ubuntu"
    hostname = aws_instance.psql.private_ip
    port     = 22
    tags     = merge(local.required_tags, var.resource_tags)
  }
}

resource "sdm_role_grant" "admin_grant_psql_ssh" {
  role_id     = sdm_role.admins.id
  resource_id = sdm_resource.psql_ssh.id
}

# SDM Roles

resource "sdm_role" "admins" {
  name = "sdm-admin-role"
}
resource "sdm_account" "admin_users" {
  count = length(var.admin_users)
  user {
    first_name = split("@", var.admin_users[count.index])[0]
    last_name  = "Onboarding"
    email      = var.admin_users[count.index]
  }
}
resource "sdm_account_attachment" "admin_attachment" {
  count      = length(var.admin_users)
  account_id = sdm_account.admin_users[count.index].id
  role_id    = sdm_role.admins.id
}

# SDM Public Key

data "sdm_ssh_ca_pubkey" "ssh_pubkey_query" {
}

# SSH Key

resource "tls_private_key" "terraform_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "terraform_key" {
  key_name   = var.key_name
  public_key = tls_private_key.terraform_key.public_key_openssh
}