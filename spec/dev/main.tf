module "api_gateway" {
    source = "terraform-aws-modules/apigateway-v2/aws"

    name          = "kaixo"
    description   = "gateway for mounting functions to dev.kaixo.io"
    protocol_type = "HTTP"

    hosted_zone_name = "dev.kaixo.io"
    domain_name      = "dev.kaixo.io"

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
        "subnet-0f00afaf6cb110510",
        "subnet-05129d09890492ffd"
    ]

    security_group_ids = [
        "sg-0018ec3d366c44cc1"
    ]
}