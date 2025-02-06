

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |

## Modules

No modules.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_api_gateway_authorizer_id"></a> [api\_gateway\_authorizer\_id](#input\_api\_gateway\_authorizer\_id) | api gateway authorizer id | `string` | `null` | no |
| <a name="input_api_gateway_authorizer_type"></a> [api\_gateway\_authorizer\_type](#input\_api\_gateway\_authorizer\_type) | api gateway authorizer type | `string` | `null` | no |
| <a name="input_api_gateway_enable"></a> [api\_gateway\_enable](#input\_api\_gateway\_enable) | enable api gateway | `bool` | `null` | no |
| <a name="input_api_gateway_id"></a> [api\_gateway\_id](#input\_api\_gateway\_id) | api gateway id | `string` | `null` | no |
| <a name="input_feature_name"></a> [feature\_name](#input\_feature\_name) | name of the feature | `string` | n/a | yes |
| <a name="input_hub_account"></a> [hub\_account](#input\_hub\_account) | hub account | <pre>object({<br/>        account_id = string<br/>        account_name = string<br/>        account_region = string<br/>    })</pre> | n/a | yes |
| <a name="input_lambda_region"></a> [lambda\_region](#input\_lambda\_region) | lambda region | `string` | `null` | no |
| <a name="input_prefix_resource_names_with_org"></a> [prefix\_resource\_names\_with\_org](#input\_prefix\_resource\_names\_with\_org) | prefix resource names with org | `bool` | `null` | no |
| <a name="input_prefix_resources_paths_with_org"></a> [prefix\_resources\_paths\_with\_org](#input\_prefix\_resources\_paths\_with\_org) | prefix resources paths with org | `bool` | `null` | no |
| <a name="input_private_subnet_ids"></a> [private\_subnet\_ids](#input\_private\_subnet\_ids) | vpc private subnet ids | `list(string)` | `null` | no |
| <a name="input_security_group_ids"></a> [security\_group\_ids](#input\_security\_group\_ids) | vpc security group ids | `list(string)` | `null` | no |
| <a name="input_ssm_region"></a> [ssm\_region](#input\_ssm\_region) | ssm region | `string` | `null` | no |
| <a name="input_substrate_name"></a> [substrate\_name](#input\_substrate\_name) | substrate name | `string` | `"agar"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_name"></a> [name](#output\_name) | n/a |
