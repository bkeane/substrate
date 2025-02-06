locals {
    organizations = toset(flatten([
        for artifact in var.artifacts:
            artifact.organization
    ]))

    image_paths = toset(flatten([
        for artifact in var.artifacts:
            artifact.image_paths
    ]))
}
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}