# Queries from SDM CA Public Key
data "sdm_ssh_ca_pubkey" "ssh_pubkey_query" {
}

# Queries latest Ubuntu AMI
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

}

resource "aws_default_vpc" "default" {
}

resource "aws_default_subnet" "default_az1" {
  availability_zone = "${var.aws_region}a"
}

# Create SDM Security Group
resource "aws_security_group" "ssh_sg" {

  name        = var.sgName
  vpc_id      = aws_default_vpc.default.id
  description = "Sec Group for SDM SSH"

}

# Adds ingress for SDM to be able to SSH
resource "aws_security_group_rule" "ingress_rules" {

  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  description       = "ssh"
  cidr_blocks       = [aws_default_vpc.default.cidr_block]
  security_group_id = aws_security_group.ssh_sg.id

}

# Adds egress for all instances to reach internet
resource "aws_security_group_rule" "egress_rules" {

  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.ssh_sg.id

}

# Create the EC2 instance
resource "aws_instance" "linux_instance" {
  ami       = data.aws_ami.ubuntu.id
  subnet_id = aws_default_subnet.default_az1.id

  vpc_security_group_ids = [aws_security_group.ssh_sg.id]
  instance_type          = "t2.micro"

  # Run commands to configure server with strongDM CA
  user_data = data.template_file.sdm_ssh_install.rendered

  tags = {
    Name = var.instanceName
  }

  volume_tags = {
    Name = var.instanceName
  }
}

# Create the SDM server resource
resource "sdm_resource" "sshCA" {
  ssh_cert {
    name     = var.instanceName
    hostname = aws_instance.linux_instance.private_dns
    username = "ubuntu"
    port     = 22
  }
}

# Renders script into base64 encoded template
data "template_file" "sdm_ssh_install" {
  template = file("${path.module}/template/sdm_ssh_install/install.sh.tpl")
  vars = {
    SSH_PUB_KEY = "${data.sdm_ssh_ca_pubkey.ssh_pubkey_query.public_key}"
  }
}