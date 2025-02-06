locals {
    hub_account = {
        account_id = "677771948337"
        account_name = "prod"
        account_region = data.aws_region.current.name
    }

    spoke_accounts = [
        {
            account_id = "831926600600"
            account_name = "dev"
            account_region = data.aws_region.current.name
        }
    ]
}

data "aws_region" "current" {}

module "monad" {
    source = "../../modules/artifacts"
    hub_account = local.hub_account

    organization = "bkeane"
    repository = "monad"
    services = [ "monad", "go", "ruby", "python", "node"  ]
}

module "substrate" {
    source = "../../modules/substrate"
    hub_account = local.hub_account
    spoke_accounts = local.spoke_accounts

    artifacts = [
        module.monad,
    ]

    features = [
        module.api,
        module.auth0,
        module.vpc,
    ]
}
