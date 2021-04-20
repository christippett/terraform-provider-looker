terraform {
  required_version = ">= 0.13"
  required_providers {
    looker = {
      source  = "dev/looker/looker"
      version = "0.0.1"
    }
  }
}

variable "looker_base_url" {
}

variable "looker_client_id" {
}

variable "looker_client_secret" {
}

provider "looker" {
  base_url      = var.looker_base_url
  client_id     = var.looker_client_id
  client_secret = var.looker_client_secret

  workspace_id = "dev"
}
