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

provider "github" {
  token = var.fine_grained_github_token
}

data "github_repository" "repository" {
  full_name = var.repository_full_name
}

resource "github_repository_environment" "repository_environment" {
  repository  = data.github_repository.repository.name
  environment = var.environment
}

module "iam_github_oidc_role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-github-oidc-role"

  name     = "${var.project_name}-${var.environment}-github-oidc-role"
  subjects = ["${data.github_repository.repository.full_name}:*"]

  policies = {
    AmazonCognitoPowerUser   = "arn:aws:iam::aws:policy/AmazonCognitoPowerUser"
    AmazonDynamoDBFullAccess = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
    AmazonS3FullAccess       = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
    AmazonSESFullAccess      = "arn:aws:iam::aws:policy/AmazonSESFullAccess"
    CloudFrontFullAccess     = "arn:aws:iam::aws:policy/CloudFrontFullAccess"
    IAMFullAccess            = "arn:aws:iam::aws:policy/IAMFullAccess"
  }
}

module "cdn" {
  source = "git::https://github.com/rfridlender/terraform-modules.git//cdn?ref=main"

  environment                  = var.environment
  force_destroy_site_s3_bucket = true
  project_name                 = var.project_name
}

data "aws_ses_email_identity" "email_identity" {
  email = var.aws_ses_email
}

module "user_pool" {
  source = "git::https://github.com/rfridlender/terraform-modules.git//user-pool?ref=main"

  aws_ses_email     = data.aws_ses_email_identity.email_identity.email
  aws_ses_email_arn = data.aws_ses_email_identity.email_identity.arn
  environment       = var.environment
  passwordless      = true
  project_name      = var.project_name
}
