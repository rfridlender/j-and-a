locals {
  variables = {
    aws_region                = var.aws_region
    aws_ses_email             = var.aws_ses_email
    fine_grained_github_token = var.fine_grained_github_token
    iam_github_oidc_role_arn  = module.iam_github_oidc_role.arn
    project_name              = var.project_name
  }
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "aws_ses_email" {
  description = "AWS SES email"
  type        = string
}

variable "fine_grained_github_token" {
  description = "Fine-grained GitHub token"
  type        = string
}

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "repository_full_name" {
  description = "Repository full name"
  type        = string
}
