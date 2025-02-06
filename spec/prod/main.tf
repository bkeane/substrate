module "api_gateway" {
    source = "terraform-aws-modules/apigateway-v2/aws"

    name          = "kaixo"
    description   = "gateway for mounting functions to prod.kaixo.io"
    protocol_type = "HTTP"

    hosted_zone_name = "prod.kaixo.io"
    domain_name      = "prod.kaixo.io"

    authorizers = {
        "auth0" = {
            name = "auth0"
            authorizer_type = "JWT"
            identity_sources = ["$request.header.Authorization"]
            jwt_configuration = {
                issuer = "https://kaixo.us.auth0.com/",
                audience = ["https://kaixo.io"]
            }
        }
    }
}

module "api" {
    source = "../../modules/feature"
    hub_account = local.hub_account

    feature_name = "api"
    api_gateway_id = module.api_gateway.api_id
    api_gateway_enable = true
}

module "auth0" {
    source = "../../modules/feature"
    hub_account = local.hub_account

    feature_name = "auth0"
    api_gateway_authorizer_id = module.api_gateway.authorizers["auth0"].id
    api_gateway_authorizer_type = "JWT"
}

module "vpc" {
    source = "../../modules/feature"
    hub_account = local.hub_account

    feature_name = "vpc"

    private_subnet_ids = [
        "subnet-0136c58f13b5f8bf9",
        "subnet-00768158825c1f939"
    ]

    security_group_ids = [
        "sg-0102ad4ccceac2613"
    ]
}

module "role" {
    source = "../../modules/role"
    hub_account = local.hub_account

    create_oidc_provider = false

    artifacts = [
        module.monad,
    ]

    bus_names = [
        module.substrate.bus_name,
    ]
}

output "workflow" {
    value = module.role.github_workflow_yaml
}