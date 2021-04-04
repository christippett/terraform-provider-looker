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
