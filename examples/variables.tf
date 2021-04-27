variable "looker" {
  type = object({
    base_url      = string
    client_id     = string
    client_secret = string
  })
}

variable "gitlab" {
  type = object({
    base_url     = string
    token        = string
    namespace_id = number
  })
}
