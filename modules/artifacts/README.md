

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |

## Modules

No modules.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_hub_account"></a> [hub\_account](#input\_hub\_account) | The hub account | <pre>object({<br/>    account_id = string<br/>    account_name = string<br/>    account_region = string<br/>  })</pre> | n/a | yes |
| <a name="input_mutable"></a> [mutable](#input\_mutable) | The mutability of the ecr repository | `bool` | `true` | no |
| <a name="input_organization"></a> [organization](#input\_organization) | The github organization | `string` | n/a | yes |
| <a name="input_repository"></a> [repository](#input\_repository) | The github repository name | `string` | n/a | yes |
| <a name="input_scan_on_push"></a> [scan\_on\_push](#input\_scan\_on\_push) | The scan on push of the ecr repository | `bool` | `false` | no |
| <a name="input_services"></a> [services](#input\_services) | The service names within the github repository | `set(string)` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_image_paths"></a> [image\_paths](#output\_image\_paths) | n/a |
| <a name="output_organization"></a> [organization](#output\_organization) | n/a |
| <a name="output_registry_id"></a> [registry\_id](#output\_registry\_id) | n/a |
