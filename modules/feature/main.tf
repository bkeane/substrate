locals {
    ssm_path = "/substrate/${var.substrate_name}/${var.feature_name}"

    collector = [
        var.lambda_region != null ? "SUBSTRATE_LAMBDA_REGION=${var.lambda_region}" : "",
        var.ssm_region != null ? "SUBSTRATE_SSM_REGION=${var.ssm_region}" : "",
        var.prefix_resources_paths_with_org != null ? "SUBSTRATE_PREFIX_PATHS_WITH_ORG=${var.prefix_resources_paths_with_org}" : "",
        var.prefix_resource_names_with_org != null ? "SUBSTRATE_PREFIX_NAMES_WITH_ORG=${var.prefix_resource_names_with_org}" : "",
        var.api_gateway_id != null ? "SUBSTRATE_APIGATEWAY_ID=${var.api_gateway_id}" : "",
        var.api_gateway_enable != null ? "SUBSTRATE_APIGATEWAY_ENABLE=${var.api_gateway_enable}" : "",
        var.api_gateway_authorizer_id != null ? "SUBSTRATE_APIGATEWAY_AUTHORIZER_ID=${var.api_gateway_authorizer_id}" : "",
        var.api_gateway_authorizer_type != null ? "SUBSTRATE_APIGATEWAY_AUTH_TYPE=${var.api_gateway_authorizer_type}" : "",
        var.private_subnet_ids != null ? "SUBSTRATE_VPC_SUBNET_IDS=${join(",", var.private_subnet_ids)}" : "",
        var.security_group_ids != null ? "SUBSTRATE_VPC_SECURITY_GROUP_IDS=${join(",", var.security_group_ids)}" : "",
    ]
}

resource "aws_ssm_parameter" "feature" {
    name = local.ssm_path
    type = "SecureString"
    value = join("\n", compact(local.collector))
}

output "name" {
    value = var.feature_name
}