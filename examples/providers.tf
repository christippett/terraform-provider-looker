terraform {
  required_version = ">= 0.13"
  required_providers {
    looker = {
      source  = "dev/looker/looker"
      version = "0.0.1"
    }
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = ">= 3.6"
    }
  }
}

provider "looker" {
  base_url      = var.looker.base_url
  client_id     = var.looker.client_id
  client_secret = var.looker.client_secret

  workspace_id = "dev"
}

provider "gitlab" {
  base_url = var.gitlab.base_url
  token    = var.gitlab.token
}
