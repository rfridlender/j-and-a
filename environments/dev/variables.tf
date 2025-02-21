locals {
  variables = {
    AWS_REGION               = var.AWS_REGION
    AWS_SES_EMAIL            = var.AWS_SES_EMAIL
    GITHUB_TOKEN             = var.GITHUB_TOKEN
    IAM_GITHUB_OIDC_ROLE_ARN = module.iam_github_oidc_role.arn
    PROJECT_NAME             = var.PROJECT_NAME
    REPOSITORY_FULL_NAME     = var.REPOSITORY_FULL_NAME
  }
}

variable "AWS_REGION" {
  description = "AWS region"
  type        = string
}

variable "AWS_SES_EMAIL" {
  description = "AWS SES email"
  type        = string
}

variable "GITHUB_TOKEN" {
  description = "GitHub token"
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
