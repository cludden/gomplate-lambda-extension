# Example: gomplate-lambda-extension
This terraform module deploys a sample Lambda function along with this extension layer installed, which it then uses to retrieve sample credentials from the SSM Parameter Store at cold start.

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_archive"></a> [archive](#requirement\_archive) | 2.2.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | 3.72.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_archive"></a> [archive](#provider\_archive) | 2.2.0 |
| <a name="provider_aws"></a> [aws](#provider\_aws) | 3.72.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_cloudwatch_log_group.example](https://registry.terraform.io/providers/hashicorp/aws/3.72.0/docs/resources/cloudwatch_log_group) | resource |
| [aws_iam_role.example](https://registry.terraform.io/providers/hashicorp/aws/3.72.0/docs/resources/iam_role) | resource |
| [aws_lambda_function.example](https://registry.terraform.io/providers/hashicorp/aws/3.72.0/docs/resources/lambda_function) | resource |
| [aws_ssm_parameter.example](https://registry.terraform.io/providers/hashicorp/aws/3.72.0/docs/resources/ssm_parameter) | resource |
| [archive_file.example](https://registry.terraform.io/providers/hashicorp/archive/2.2.0/docs/data-sources/file) | data source |
| [aws_caller_identity.current](https://registry.terraform.io/providers/hashicorp/aws/3.72.0/docs/data-sources/caller_identity) | data source |
| [aws_iam_policy_document.example](https://registry.terraform.io/providers/hashicorp/aws/3.72.0/docs/data-sources/iam_policy_document) | data source |
| [aws_iam_policy_document.trust](https://registry.terraform.io/providers/hashicorp/aws/3.72.0/docs/data-sources/iam_policy_document) | data source |
| [aws_lambda_invocation.example](https://registry.terraform.io/providers/hashicorp/aws/3.72.0/docs/data-sources/lambda_invocation) | data source |
| [aws_partition.current](https://registry.terraform.io/providers/hashicorp/aws/3.72.0/docs/data-sources/partition) | data source |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/3.72.0/docs/data-sources/region) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_bar"></a> [bar](#input\_bar) | sample bar credential | `list(string)` | <pre>[<br>  "a",<br>  "b",<br>  "c"<br>]</pre> | no |
| <a name="input_foo"></a> [foo](#input\_foo) | sample foo credential | `string` | `"secret"` | no |
| <a name="input_name"></a> [name](#input\_name) | example component name | `string` | `"gomplate-lambda-extension-example"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_result"></a> [result](#output\_result) | lambda function config |
<!-- END_TF_DOCS -->