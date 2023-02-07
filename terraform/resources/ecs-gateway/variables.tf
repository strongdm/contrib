#SDM VARIABLES
variable "sdm_access_key" {
  type        = string
  description = "SDM Access Key"
}

variable "sdm_secret_key" {
  type        = string
  description = "SDM Secret Key"
}

#AWS Varibales
variable "region" {
  type        = string
  description = "Region for your AWS Infrastructure"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID you'd like the ECS Fargate Instance in."
}

variable "public_subnet_id" {
  type        = string
  description = "Public Subnet for the ECS Fargate Instance. Needs to be same AZ as Private Subnet"
}

variable "private_subnet_id" {
  type        = string
  description = "Private Subnet for the ECS Fargate Instance. Needs to be same AZ as Public Subnet"
}

variable "sg_sdm_ingress_rules" {
  description = "Block for security group rules. Port 5000 needs to be open by default."

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
      description = "Ingress for StrongDM Gateways"
    }
  ]
}