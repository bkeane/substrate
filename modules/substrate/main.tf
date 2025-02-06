locals {
    hub_account_id = var.hub_account.account_id != null ? var.hub_account.account_id : data.aws_caller_identity.current.account_id
    hub_account_name = var.hub_account.account_name

    is_hub = local.hub_account_id == data.aws_caller_identity.current.account_id
    is_spoke = !local.is_hub

    resource_name = var.substrate_name
    resource_path = "/substrate/${var.substrate_name}"

    features = toset(flatten([ for feature in var.features: feature.name ]))
    artifacts = toset(flatten([ for artifact in var.artifacts: artifact.image_paths ]))

    source = lookup(local.account_map, data.aws_caller_identity.current.account_id)
    destinations = local.is_hub ? [ for id, name in local.account_map: name ] : [ local.source ]

    account_map = merge(local.spoke_map, local.hub_map)
    spoke_map = { for spoke in var.spoke_accounts: spoke.account_id => spoke.account_name }
    hub_map = { (local.hub_account_id) = local.hub_account_name }
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

resource "aws_ssm_parameter" "substrate" {
    name = local.resource_path
    type = "SecureString"
    value = <<EOF
SUBSTRATE_NAME=${var.substrate_name}
SUBSTRATE_SOURCE=${local.source}
SUBSTRATE_DESTINATIONS=${join(",", local.destinations)}
SUBSTRATE_FEATURES=${join(",", local.features)}
SUBSTRATE_ECR_REGISTRY_ID=${local.hub_account_id}
SUBSTRATE_ECR_REGISTRY_REGION=${data.aws_region.current.name}
SUBSTRATE_EVENTBRIDGE_BUS_NAME=${aws_cloudwatch_event_bus.backplane.name}
SUBSTRATE_EVENTBRIDGE_ENABLE=true
SUBSTRATE_APIGATEWAY_ENABLE=false
EOF
}

output "bus_name" {
    value = aws_cloudwatch_event_bus.backplane.name
}
