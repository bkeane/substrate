variable "substrate_name" {
  type = string
  description = "The name of the substrate"
  default = "agar"
}

variable "hub_account" {
  type = object({
    account_id = string
    account_name = string
    account_region = string
  })
  description = "The Hub account, defaults to the caller account"
}

variable "spoke_accounts" {
  type = set(object({
    account_id = string
    account_name = string
    account_region = string
  }))
  description = "The IDs of the Spoke accounts (account_name => account_id)"
  default = []
}

variable "artifacts" {
  type = set(object({
    organization = string
    registry_id = string
    image_paths = set(string)
  }))
  description = "Artifacts to be made distributable by the substrate"
  default = []
}

variable "features" {
  type = set(object({
    name = string
  }))
  description = "Features to be made usable via the substrate"
  default = []
}