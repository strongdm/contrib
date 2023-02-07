provider "aws" {
  region = var.region
}

provider "sdm" {
  api_access_key = var.sdm_access_key
  api_secret_key = var.sdm_secret_key
}

resource "random_id" "id" {
  byte_length = 6
}
