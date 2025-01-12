resource "github_actions_environment_variable" "environment_variables" {
  for_each = {
    "AWS_REGION"               = var.aws_region
    "AWS_SES_EMAIL"            = var.aws_ses_email
    "AWS_SES_EMAIL_ARN"        = var.aws_ses_email_arn
    "ENVIRONMENT"              = var.environment
    "IAM_GITHUB_OIDC_ROLE_ARN" = module.iam_github_oidc_role.arn
    "PROJECT_NAME"             = var.project_name
  }

  repository    = var.project_name
  environment   = github_repository_environment.repository_environment.environment
  variable_name = each.key
  value         = each.value
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "aws_ses_email" {
  description = "AWS SES email"
  type        = string
}

variable "aws_ses_email_arn" {
  description = "AWS SES email ARN"
  type        = string
}

variable "environment" {
  description = "Environment"
  type        = string
}

variable "github_account_name" {
  description = "GitHub account name"
  type        = string
}

variable "project_name" {
  description = "Project name"
  type        = string
}
