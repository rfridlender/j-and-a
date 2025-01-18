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
    AmazonAPIGatewayAdministrator = "arn:aws:iam::aws:policy/AmazonAPIGatewayAdministrator"
    AmazonCognitoPowerUser        = "arn:aws:iam::aws:policy/AmazonCognitoPowerUser"
    AmazonDynamoDBFullAccess      = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
    AmazonS3FullAccess            = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
    AmazonSESFullAccess           = "arn:aws:iam::aws:policy/AmazonSESFullAccess"
    AWSLambda_FullAccess          = "arn:aws:iam::aws:policy/AWSLambda_FullAccess"
    CloudFrontFullAccess          = "arn:aws:iam::aws:policy/CloudFrontFullAccess"
    CloudWatchFullAccessV2        = "arn:aws:iam::aws:policy/CloudWatchFullAccessV2"
    IAMFullAccess                 = "arn:aws:iam::aws:policy/IAMFullAccess"
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

locals {
  dynamodb_index_name = "${var.project_name}-${var.environment}-EntityType-SK-index"
}

module "dynamodb_table" {
  source = "terraform-aws-modules/dynamodb-table/aws"

  name      = "${var.project_name}-${var.environment}-PK-SK-table"
  hash_key  = "PK"
  range_key = "SK"

  global_secondary_indexes = [
    {
      name            = local.dynamodb_index_name
      hash_key        = "EntityType"
      range_key       = "SK"
      projection_type = "ALL"
    }
  ]

  attributes = [
    {
      name = "PK"
      type = "S"
    },
    {
      name = "SK"
      type = "S"
    },
    {
      name = "EntityType"
      type = "S"
    },
  ]
}

module "artifact_store" {
  source = "terraform-aws-modules/s3-bucket/aws"

  bucket = "${var.project_name}-${var.environment}-artifact-store"
  acl    = "private"

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true

  control_object_ownership = true
  object_ownership         = "ObjectWriter"

  versioning = {
    enabled = true
  }
}

module "api_gateway" {
  source = "git::https://github.com/rfridlender/terraform-modules.git//api-gateway?ref=main"

  aws_region   = var.aws_region
  environment  = var.environment
  project_name = var.project_name
  routes = {
    "GET /jobs/{jobId}/logs"            = module.function_log.lambda_function_arn
    "GET /jobs/{jobId}/logs/{logId}"    = module.function_log.lambda_function_arn
    "PUT /jobs/{jobId}/logs/{logId}"    = module.function_log.lambda_function_arn
    "DELETE /jobs/{jobId}/logs/{logId}" = module.function_log.lambda_function_arn
  }
  user_pool_id         = module.user_pool.user_pool_id
  user_pool_client_ids = [module.user_pool.user_pool_client_id]
}

module "function_iam_policy" {
  source = "terraform-aws-modules/iam/aws//modules/iam-policy"

  name = "${var.project_name}-${var.environment}-function-iam-policy"
  path = "/"

  policy = templatefile(
    "${path.module}/tpls/function-iam-policy.json",
    { dynamodb_table_arn = module.dynamodb_table.dynamodb_table_arn },
  )
}

module "function_log" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "${var.project_name}-${var.environment}-function-log"
  runtime       = "provided.al2023"
  handler       = "bootstrap"
  architectures = ["arm64"]
  publish       = true

  source_path = "../../cmd/log/bootstrap"

  store_on_s3 = true
  s3_bucket   = module.artifact_store.s3_bucket_id

  allowed_triggers = {
    api_gateway = {
      service    = "apigateway"
      source_arn = "${module.api_gateway.stage_execution_arn}/*/*"
    }
  }

  attach_policy = true
  policy        = module.function_iam_policy.arn

  environment_variables = {
    DYNAMO_DB_TABLE_NAME = module.dynamodb_table.dynamodb_table_id
    DYNAMO_DB_INDEX_NAME = local.dynamodb_index_name
  }
}

