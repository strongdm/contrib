# Provider

terraform {
  required_version = ">= 0.15.0"
  required_providers {
    aws = ">= 3.0.0"
    sdm = {
      source  = "strongdm/sdm"
      version = ">= 1.0.12"
    }
  }
}

provider "aws" {
  region  = var.region
}

provider "sdm" {
  api_access_key = var.sdm_access_key
  api_secret_key = var.sdm_secret_key
}
