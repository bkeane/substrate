locals {
    is_hub = var.hub_account.account_id == data.aws_caller_identity.current.account_id
    mutability = var.mutable ? "MUTABLE" : "IMMUTABLE"
    image_paths = toset(flatten([
        for service in var.services:
            "${var.organization}/${var.repository}/${service}"
    ]))
}

data "aws_caller_identity" "current" {}

resource "aws_ecr_repository" "service" {
    for_each = local.is_hub ? local.image_paths : []
    name = each.value
    image_tag_mutability = local.mutability
    image_scanning_configuration {
        scan_on_push = var.scan_on_push
    }
}
