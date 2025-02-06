

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |

## Modules

No modules.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_artifacts"></a> [artifacts](#input\_artifacts) | Artifacts to be made distributable by the substrate | <pre>set(object({<br/>    organization = string<br/>    registry_id = string<br/>    image_paths = set(string)<br/>  }))</pre> | `[]` | no |
| <a name="input_features"></a> [features](#input\_features) | Features to be made usable via the substrate | <pre>set(object({<br/>    name = string<br/>  }))</pre> | `[]` | no |
| <a name="input_hub_account"></a> [hub\_account](#input\_hub\_account) | The Hub account, defaults to the caller account | <pre>object({<br/>    account_id = string<br/>    account_name = string<br/>    account_region = string<br/>  })</pre> | n/a | yes |
| <a name="input_spoke_accounts"></a> [spoke\_accounts](#input\_spoke\_accounts) | The IDs of the Spoke accounts (account\_name => account\_id) | <pre>set(object({<br/>    account_id = string<br/>    account_name = string<br/>    account_region = string<br/>  }))</pre> | `[]` | no |
| <a name="input_substrate_name"></a> [substrate\_name](#input\_substrate\_name) | The name of the substrate | `string` | `"agar"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_bus_name"></a> [bus\_name](#output\_bus\_name) | n/a |
