terraform {
  backend "s3" {
    bucket         = "rfridlender-terraform-state"
    key            = "terraform-only-stack/environments/dev/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "rfridlender-terraform-state"
  }
}
