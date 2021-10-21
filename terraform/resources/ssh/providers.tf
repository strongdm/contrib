terraform {
  required_providers {
    sdm = {
      source  = "strongdm/sdm"
      version = ">= 1.0.12"
    }
  }
}

provider "aws" {
  region = var.aws_region
}