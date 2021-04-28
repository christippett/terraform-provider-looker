data "looker_session" "test" {}

resource "looker_project" "test" {
  name = "test_project"
}

resource "looker_user" "test" {
  first_name = "Lloyd"
  last_name  = "Tabb"

  credentials_email {
    email = "lloyd.tabb@example.com"
  }
}

resource "looker_git_deploy_key" "test" {
  project = looker_project.test.id

  depends_on = [looker_project.test]
}

resource "gitlab_project" "looker" {
  name             = looker_project.test.name
  description      = "Looker test repository created by Terraform"
  visibility_level = "private"
  namespace_id     = var.gitlab.namespace_id

  initialize_with_readme = true
}

resource "gitlab_deploy_key" "test" {
  project  = gitlab_project.looker.path_with_namespace
  title    = "Looker git deploy key"
  key      = looker_git_deploy_key.test.public_key
  can_push = true
}

resource "looker_project_git_repo" "test" {
  project                            = looker_project.test.id
  git_remote_url                     = gitlab_project.looker.ssh_url_to_repo
  git_service_name                   = "gitlab"
  git_application_server_http_port   = 443
  git_application_server_http_scheme = "https"

  depends_on = [gitlab_deploy_key.test]
}

output "looker_session" {
  value = data.looker_session.test
}
