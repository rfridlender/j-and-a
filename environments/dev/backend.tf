terraform {
  backend "s3" {
    bucket         = "rfridlender-terraform-state"
    key            = "j-and-a/environments/dev/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "rfridlender-terraform-state"
  }
}
