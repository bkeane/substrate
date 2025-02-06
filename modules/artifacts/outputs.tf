output "organization" {
    value = var.organization
}

output "registry_id" {
    value = var.hub_account.account_id
}

output "image_paths" {
    value = local.image_paths
}