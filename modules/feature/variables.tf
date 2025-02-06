variable "substrate_name" {
    description = "substrate name"
    type = string
    default = "agar"
}

variable "hub_account" {
    description = "hub account"
    type = object({
        account_id = string
        account_name = string
        account_region = string
    })
}

variable "feature_name" {
    description = "name of the feature"
    type = string
}

variable "lambda_region" {
    description = "lambda region"
    type = string
    default = null
}

variable "ssm_region" {
    description = "ssm region"
    type = string
    default = null
}

variable "prefix_resources_paths_with_org" {
    description = "prefix resources paths with org"
    type = bool
    default = null
}

variable "prefix_resource_names_with_org" {
    description = "prefix resource names with org"
    type = bool
    default = null
}

variable "api_gateway_id" {
    description = "api gateway id"
    type = string
    default = null
}

variable "api_gateway_enable" {
    description = "enable api gateway"
    type = bool
    default = null
}

variable "api_gateway_authorizer_id" {
    description = "api gateway authorizer id"
    type = string
    default = null
}

variable "api_gateway_authorizer_type" {
    description = "api gateway authorizer type"
    type = string
    default = null
}

variable "private_subnet_ids" {
    description = "vpc private subnet ids"
    type = list(string)
    default = null
}

variable "security_group_ids" {
    description = "vpc security group ids"
    type = list(string)
    default = null
}