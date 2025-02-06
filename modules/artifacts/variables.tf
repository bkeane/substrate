variable "hub_account" {
  description = "The hub account"
  type = object({
    account_id = string
    account_name = string
    account_region = string
  })
}

variable "organization" {
  description = "The github organization"
  type = string
}

variable "repository" {
  description = "The github repository name"
  type = string
}

variable "services" {
  description = "The service names within the github repository"
  type = set(string)
}

variable "mutable" {
  description = "The mutability of the ecr repository"
  type = bool
  default = true
}

variable "scan_on_push" {
  description = "The scan on push of the ecr repository"
  type = bool
  default = false
}