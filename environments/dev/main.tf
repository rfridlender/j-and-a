terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }

    github = {
      source  = "integrations/github"
      version = "~> 6.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

resource "github_repository_environment" "repository_environment" {
  repository  = var.project_name
  environment = var.environment
}

module "iam_github_oidc_role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-github-oidc-role"

  name     = "${var.project_name}-${var.environment}-github-oidc-role"
  subjects = ["${var.github_account_name}/${var.project_name}:*"]

  policies = {
    AmazonCognitoPowerUser   = "arn:aws:iam::aws:policy/AmazonCognitoPowerUser"
    AmazonDynamoDBFullAccess = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
    AmazonS3FullAccess       = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
    CloudFrontFullAccess     = "arn:aws:iam::aws:policy/CloudFrontFullAccess"
    IAMFullAccess            = "arn:aws:iam::aws:policy/IAMFullAccess"
  }
}

module "cdn" {
  source = "git::https://github.com/rfridlender/terraform-modules.git//cdn?ref=main"

  environment   = var.environment
  force_destroy = true
  project_name  = var.project_name
}

module "user_pool" {
  source = "git::https://github.com/rfridlender/terraform-modules.git//user-pool?ref=main"

  aws_ses_email     = var.aws_ses_email
  aws_ses_email_arn = var.aws_ses_email_arn
  environment       = var.environment
  passwordless      = true
  project_name      = var.project_name
}
