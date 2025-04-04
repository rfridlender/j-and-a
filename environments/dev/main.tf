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
  region = var.AWS_REGION
}

provider "github" {
  token = var.FINE_GRAINED_GITHUB_TOKEN
}

locals {
  environment = basename(abspath(path.module))
}

data "github_repository" "repository" {
  full_name = var.REPOSITORY_FULL_NAME
}

resource "github_repository_environment" "repository_environment" {
  repository  = data.github_repository.repository.name
  environment = local.environment
}

resource "github_actions_environment_variable" "tf_vars_json_environment_variable" {
  repository    = var.PROJECT_NAME
  environment   = github_repository_environment.repository_environment.environment
  variable_name = "TF_VARS_JSON"
  value         = jsonencode(local.environment_variables)
}

resource "github_actions_environment_variable" "pipeline_dependent_environment_variables" {
  for_each      = local.pipeline_dependent_environment_variable_keys
  repository    = var.PROJECT_NAME
  environment   = github_repository_environment.repository_environment.environment
  variable_name = each.value
  value         = local.environment_variables[each.value]
}

module "iam_github_oidc_role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-github-oidc-role"

  name     = "${var.PROJECT_NAME}-${local.environment}-github-oidc-role"
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

  environment                  = local.environment
  force_destroy_site_s3_bucket = true
  project_name                 = var.PROJECT_NAME
}

data "aws_ses_email_identity" "email_identity" {
  email = var.AWS_SES_EMAIL
}

module "user_pool" {
  source = "git::https://github.com/rfridlender/terraform-modules.git//user-pool?ref=main"

  aws_ses_email     = data.aws_ses_email_identity.email_identity.email
  aws_ses_email_arn = data.aws_ses_email_identity.email_identity.arn
  environment       = local.environment
  passwordless      = true
  project_name      = var.PROJECT_NAME
}

locals {
  dynamodb_index_name = "${var.PROJECT_NAME}-${local.environment}-ModelType-SK-index"
}

module "dynamodb_table" {
  source = "terraform-aws-modules/dynamodb-table/aws"

  name      = "${var.PROJECT_NAME}-${local.environment}-PK-SK-table"
  hash_key  = "PK"
  range_key = "SK"

  global_secondary_indexes = [
    {
      name            = local.dynamodb_index_name
      hash_key        = "ModelType"
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
      name = "ModelType"
      type = "S"
    },
  ]
}

module "artifact_store" {
  source = "terraform-aws-modules/s3-bucket/aws"

  bucket = "${var.PROJECT_NAME}-${local.environment}-artifact-store"
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

  aws_region   = var.AWS_REGION
  environment  = local.environment
  project_name = var.PROJECT_NAME
  routes = {
    "DELETE /{PartitionType}/{PartitionId}/{SortType}"          = module.function_model.lambda_function_arn
    "DELETE /{PartitionType}/{PartitionId}/{SortType}/{SortId}" = module.function_model.lambda_function_arn
    "GET /{PartitionType}/{PartitionId}/{SortType}"             = module.function_model.lambda_function_arn
    "GET /{PartitionType}/{PartitionId}/{SortType}/{SortId}"    = module.function_model.lambda_function_arn
    "GET /{SortType}"                                           = module.function_model.lambda_function_arn
    "PUT /{PartitionType}/{PartitionId}/{SortType}"             = module.function_model.lambda_function_arn
    "PUT /{PartitionType}/{PartitionId}/{SortType}/{SortId}"    = module.function_model.lambda_function_arn
  }
  user_pool_id         = module.user_pool.user_pool_id
  user_pool_client_ids = [module.user_pool.user_pool_client_id]
}

module "function_iam_policy" {
  source = "terraform-aws-modules/iam/aws//modules/iam-policy"

  name = "${var.PROJECT_NAME}-${local.environment}-function-iam-policy"
  path = "/"

  policy = templatefile(
    "${path.module}/tpls/function-iam-policy.json",
    { dynamodb_table_arn = module.dynamodb_table.dynamodb_table_arn },
  )
}

module "function_model" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "${var.PROJECT_NAME}-${local.environment}-function-model"
  runtime       = "provided.al2023"
  handler       = "bootstrap"
  architectures = ["arm64"]
  publish       = true

  source_path = "../../cmd/function-model/bootstrap"

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
