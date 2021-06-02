data "sdm_ssh_ca_pubkey" "ssh_pubkey_query" {
  # returns the SDM-generated CA for your org
}

data "aws_ami" "ubuntu" {
    most_recent = true
    owners = ["099720109477"] # Canonical
}

resource "aws_default_vpc" "default" {
}

resource "aws_default_subnet" "default_az1" {
  availability_zone = "${var.aws_region}a"
}

# Create the EC2 instance
resource "aws_instance" "linux_instance" {
  ami             = data.aws_ami.ubuntu.id
  subnet_id       = aws_default_subnet.default_az1.id

  # security_groups = var.securityGroups 
  instance_type   = "t2.micro"

  # Run commands to configure server with strongDM CA
  user_data = <<USERDATA
#!/bin/bash -xe
sudo su -
sudo echo '${data.sdm_ssh_ca_pubkey.ssh_pubkey_query.public_key}' >> /etc/ssh/sdm_ca.pub
chmod 0600 /etc/ssh/sdm_ca.pub
echo 'TrustedUserCAKeys /etc/ssh/sdm_ca.pub' | tee -a /etc/ssh/sshd_config
service ssh restart
USERDATA
  
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
    hostname = aws_instance.linux_instance.public_dns
    username = "ubuntu"
    port     = 22
  }
}