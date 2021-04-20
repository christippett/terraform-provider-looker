resource "looker_project" "test" {
  name = "test_project"
}

resource "looker_user" "test" {
  first_name = "John"
  last_name  = "Smith"

  credentials_email {
    email = "john.smith@example.com"
  }
}


data "looker_session" "test" {}

data "http" "test" {
  # url =
  # "https://servian.eu.looker.com:19999/api/3.1/projects/test_project/git/deploy_key"
  url = "https://servian.eu.looker.com:19999/api/3.1/session"

  # Optional request headers
  request_headers = {
    Authorization = "Bearer ${data.looker_session.test.access_token}"
  }

  depends_on = [
    looker_project.test,
    data.looker_session.test
  ]
}

output "test_output" {
  value = data.looker_session.test
}

output "session" {
  value = data.http.test
}

# output "git_deploy_key" {
#   # value = looker_project.test.git_deploy_key
#   value = data.http.test
# }
