variable "bar" {
  description = "sample bar credential"
  type        = list(string)
  sensitive   = true
  default     = ["a", "b", "c"]
}

variable "foo" {
  description = "sample foo credential"
  type        = string
  sensitive   = true
  default     = "secret"
}

variable "name" {
  description = "example component name"
  type        = string
  default     = "gomplate-lambda-extension-example"
}

variable "release" {
  description = "extension release name"
  type        = string
  default     = "0.1.1"
}
