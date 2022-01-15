output "result" {
  description = "lambda function config"
  value       = jsondecode(data.aws_lambda_invocation.example.result)
}
