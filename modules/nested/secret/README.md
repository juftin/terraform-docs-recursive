<!-- BEGIN_TF_DOCS -->
# secret

This module creates a secret in AWS Secrets Manager.
This is also a nested module, which will still be found.

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | ~> 5.0 ~> 5.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_secretsmanager_secret.secret](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/secretsmanager_secret) | resource |
| [aws_secretsmanager_secret_version.secret_version](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/secretsmanager_secret_version) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_region"></a> [region](#input\_region) | n/a | `string` | `"us-east-1"` | no |
| <a name="input_secret_name"></a> [secret\_name](#input\_secret\_name) | n/a | `string` | `"my-test-bucket"` | no |
| <a name="input_secret_value"></a> [secret\_value](#input\_secret\_value) | n/a | `string` | `"my-test-bucket"` | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->
