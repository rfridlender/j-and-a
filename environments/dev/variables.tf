locals {
  environment_variables = tomap({
    AWS_REGION                = var.AWS_REGION
    AWS_SES_EMAIL             = var.AWS_SES_EMAIL
    FINE_GRAINED_GITHUB_TOKEN = var.FINE_GRAINED_GITHUB_TOKEN
    IAM_GITHUB_OIDC_ROLE_ARN  = module.iam_github_oidc_role.arn
    PROJECT_NAME              = var.PROJECT_NAME
    REPOSITORY_FULL_NAME      = var.REPOSITORY_FULL_NAME
  })

  pipeline_dependent_environment_variable_keys = toset(["AWS_REGION", "FINE_GRAINED_GITHUB_TOKEN", "IAM_GITHUB_OIDC_ROLE_ARN"])
}

variable "AWS_REGION" {
  description = "AWS region"
  type        = string
}

variable "AWS_SES_EMAIL" {
  description = "AWS SES email"
  type        = string
}

variable "FINE_GRAINED_GITHUB_TOKEN" {
  description = "Fine-grained GitHub token"
  type        = string
}

variable "PROJECT_NAME" {
  description = "Project name"
  type        = string
}

variable "REPOSITORY_FULL_NAME" {
  description = "Repository full name"
  type        = string
}
