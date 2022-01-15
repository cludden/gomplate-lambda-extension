# gomplate-lambda-extension
a [Gomplate](https://github.com/hairyhenderson/gomplate) extension for [AWS Lambda](https://aws.amazon.com/blogs/compute/introducing-aws-lambda-extensions-in-preview/)

## Getting Started

The [example](./example) directory provides a sample Terraform module that can be used to provision a Lambda function with this extension enabled.

*Prereqs:*
- [Terraform v1+](https://www.terraform.io/downloads)
- [An AWS Account](https://aws.amazon.com)

```shell
cd example
terraform init
terraform apply
```

## Installation

**Lambda layer**
For Lambda functions that use .zip deployments, include the following ARN as a layer in your Lambda function:
```
arn:aws:lambda:us-west-2:010013098410:layer:gomplate-lambda-extension:4
```

or publish your own layer:
```shell
just publish
```

**Embed within Lambda container image**  
For Lambda functions that use container images, install the extension as part of the container build:
```dockerfile
ADD https://github.com/cludden/gomplate-lambda-extension/releases/latest/download/gomplate-lambda-extension /opt/extensions/gomplate-lambda-extension
```


## Configuration
The extension is configured via [Lambda environment variables](https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html).

| Name | Description | Required |
| :--- | :--- | :---: |
| `GOMPLATE_DATASOURCE_{foo}` | defines a named datasource | |
| `GOMPLATE_INPUT` | input template | `true` |
| `GOMPLATE_OUTPUT` | output file path (e.g. `/tmp/config.json`) | `true` |

## License
Licensed under the [MIT License](LICENSE.md)  
Copyright (c) 2022 Chris Ludden