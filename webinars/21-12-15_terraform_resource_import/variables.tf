variable "region" {
  description = "AWS Region To Deploy"
}

variable "sdm_access_key" {
  description = "SDM API Access Key"
}

variable "sdm_secret_key" {
  description = "SDM API Secret Key"
}

variable "project_name" {
  type    = string
  default = "sdm-dev"
}

variable "main_vpc_cidr" {
  default     = "10.0.0.0/16"
  description = "VPC CIDR Range"
}

variable "public_subnets" {
  default     = "10.0.1.0/24"
  description = "Public Subnet CIDR Range"
}

variable "private_subnets" {
  default     = "10.0.2.0/24"
  description = "Private Subnet CIDR Range"
}

variable "resource_tags" {
  type = map(string)
  default = {
    Terraform = "true"
    Creator   = "John T."
  }
}

variable "admin_users" {
  type        = list(string)
  default     = []
  description = "A list of email addresses that will be granted access to all resources."
}

variable "psql_sg_ingress_rules" {

  type = list(object({

    from_port   = number
    to_port     = number
    protocol    = string
    cidr_block  = string
    description = string

  }))

  default = [
    {
      from_port   = 5432
      to_port     = 5432
      protocol    = "tcp"
      cidr_block  = "10.0.2.0/24"
      description = "psql"
    },
    {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_block  = "10.0.2.0/24"
      description = "ssh"
    }
  ]

}

variable "gateway_sg_ingress_rules" {

  type = list(object({

    from_port   = number
    to_port     = number
    protocol    = string
    cidr_block  = string
    description = string

  }))

  default = [
    {
      from_port   = 5000
      to_port     = 5000
      protocol    = "tcp"
      cidr_block  = "0.0.0.0/0"
      description = "sdm"
    },
    {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_block  = "0.0.0.0/0"
      description = "ssh"
    }
  ]

}

variable "vault_sg_ingress_rules" {

  type = list(object({

    from_port   = number
    to_port     = number
    protocol    = string
    cidr_block  = string
    description = string

  }))

  default = [
    {
      from_port   = 8200
      to_port     = 8201
      protocol    = "tcp"
      cidr_block  = "10.0.2.0/24"
      description = "vault"
    },
    {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_block  = "10.0.2.0/24"
      description = "ssh"
    },
    {
      from_port   = 8500
      to_port     = 8501
      protocol    = "tcp"
      cidr_block  = "10.0.2.0/24"
      description = "consul"
    }
  ]

}

variable "key_name" {
  default = "jturner-terraform-key"
}