################################################################################
## Terraform
################################################################################

terraform {
  required_providers {
    archive = {
      source  = "hashicorp/archive"
      version = "2.2.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "3.72.0"
    }
    get = {
      source  = "cludden/get"
      version = "0.1.0"
    }
  }
}

################################################################################
## Data Sources
################################################################################

# create zip artifact with sample function
data "archive_file" "example" {
  type             = "zip"
  source_file      = "${path.root}/index.mjs"
  output_file_mode = "0666"
  output_path      = "${path.root}/index.mjs.zip"
}

# lookup aws session
data "aws_caller_identity" "current" {}

# lookup aws partition info
data "aws_partition" "current" {}

# lookup aws region info
data "aws_region" "current" {}

################################################################################
## Resources
################################################################################

# invoke lambda function
data "aws_lambda_invocation" "example" {
  function_name = aws_lambda_function.example.function_name
  input         = jsonencode({})
}

# provision lambda function
resource "aws_lambda_function" "example" {
  filename         = data.archive_file.example.output_path
  function_name    = var.name
  handler          = "index.handler"
  layers           = [aws_lambda_layer_version.gomplate.arn]
  role             = aws_iam_role.example.arn
  runtime          = "nodejs14.x"
  source_code_hash = filebase64sha256(data.archive_file.example.output_path)
  timeout          = 30

  environment {
    variables = {
      GOMPLATE_DATASOURCE_ssm = "aws+smp:///${var.name}/"
      GOMPLATE_INPUT          = <<-JSON
        {
          "foo": "{{ (ds "ssm" "foo").Value }}",
          "bar": {{ (ds "ssm" "bar").Value | strings.Split "," | data.ToJSON }}
        }
      JSON
      GOMPLATE_OUTPUT         = "/tmp/config.json"
    }
  }

  depends_on = [
    aws_cloudwatch_log_group.example
  ]
}

# provision extension lambda layer
resource "aws_lambda_layer_version" "gomplate" {
  filename            = get_artifact.extension.dest
  layer_name          = "${var.name}-gomplate"
  compatible_runtimes = ["go1.x"]
}

# download extension
resource "get_artifact" "extension" {
  url      = "https://github.com/cludden/gomplate-lambda-extension/releases/download/v${var.release}/gomplate-lambda-extension_${var.release}_linux_amd64.zip"
  checksum = "file:https://github.com/cludden/gomplate-lambda-extension/releases/download/v${var.release}/checksums.txt"
  dest     = "gomplate-lambda-extension_${var.release}_linux_amd64.zip"
  mode     = "file"
  archive  = false
  workdir  = abspath(path.root)
}

##############################
## CloudWatch
##############################

# provision log group for function logs
resource "aws_cloudwatch_log_group" "example" {
  name              = "/aws/lambda/${var.name}"
  retention_in_days = 7
}

##############################
## IAM
##############################

# provision execution role
resource "aws_iam_role" "example" {
  name               = var.name
  assume_role_policy = data.aws_iam_policy_document.trust.json

  inline_policy {
    name   = "inline"
    policy = data.aws_iam_policy_document.example.json
  }
}

# define execution role trust policy
data "aws_iam_policy_document" "trust" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

# define execution role policy statements
data "aws_iam_policy_document" "example" {
  # allow function to push logs to cloudwatch
  statement {
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    effect    = "Allow"
    resources = ["${aws_cloudwatch_log_group.example.arn}*"]
  }

  # allow function to read ssm parameters
  statement {
    actions   = ["ssm:GetParameter"]
    effect    = "Allow"
    resources = ["arn:${data.aws_partition.current.partition}:ssm:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:parameter/${var.name}/*"]
  }
}

##############################
## SSM
##############################

# provision sample secrets as ssm parameters
resource "aws_ssm_parameter" "example" {
  for_each = {
    foo = var.foo
    bar = join(",", var.bar)
  }
  name  = "/${var.name}/${each.key}"
  type  = "SecureString"
  value = each.value
}
