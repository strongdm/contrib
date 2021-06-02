terraform {
  required_providers {
    sdm = {
      source = "strongdm/sdm"
    }
  }
}

provider "aws" {
    region = var.aws_region
}
