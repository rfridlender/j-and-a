output "cloudfront_distribution_id" {
  value       = module.cdn.cloudfront_distribution_id
  description = "CloudFront distribution ID"
}

output "s3_bucket_name" {
  value       = module.cdn.s3_bucket_name
  description = "S3 bucket name"
}

output "vite_user_pool_id" {
  value       = module.user_pool.vite_user_pool_id
  description = "Vite user pool ID"
}

output "vite_user_pool_client_id" {
  value       = module.user_pool.vite_user_pool_client_id
  description = "Vite pool client ID"
}
