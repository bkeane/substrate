variable "name" {
  type = string
  description = "name for resources"
  default = "agar"
}

variable "hub_account" {
  type = object({
    account_id = string
    account_name = string
    account_region = string
  })
  description = "The Hub account"
}

variable "artifacts" {
  type = set(object({
    organization = string
    registry_id = string
    image_paths = set(string)
  }))
  description = "The artifacts to allow"
}

variable "bus_names" {
  type = set(string)
  description = "The name of the event buses to allow"
}

variable "create_oidc_provider" {
  type = bool
  description = "Whether to create the OIDC provider or lookup existing"
  default = false
}