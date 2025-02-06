

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |

## Modules

No modules.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_artifacts"></a> [artifacts](#input\_artifacts) | The artifacts to allow | <pre>set(object({<br/>    organization = string<br/>    registry_id = string<br/>    image_paths = set(string)<br/>  }))</pre> | n/a | yes |
| <a name="input_bus_names"></a> [bus\_names](#input\_bus\_names) | The name of the event buses to allow | `set(string)` | n/a | yes |
| <a name="input_create_oidc_provider"></a> [create\_oidc\_provider](#input\_create\_oidc\_provider) | Whether to create the OIDC provider or lookup existing | `bool` | `false` | no |
| <a name="input_hub_account"></a> [hub\_account](#input\_hub\_account) | The Hub account | <pre>object({<br/>    account_id = string<br/>    account_name = string<br/>    account_region = string<br/>  })</pre> | n/a | yes |
| <a name="input_name"></a> [name](#input\_name) | name for resources | `string` | `"agar"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_github_workflow_yaml"></a> [github\_workflow\_yaml](#output\_github\_workflow\_yaml) | n/a |
| <a name="output_role_arn"></a> [role\_arn](#output\_role\_arn) | n/a |
