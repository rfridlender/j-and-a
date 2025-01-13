output "cloudfront_distribution_id" {
  value       = module.cdn.cloudfront_distribution_id
  description = "CloudFront distribution ID"
}

output "site_s3_bucket_name" {
  value       = module.cdn.site_s3_bucket_name
  description = "Site S3 bucket name"
}

output "user_pool_id" {
  value       = module.user_pool.vite_user_pool_id
  description = "User pool ID"
}

output "user_pool_client_id" {
  value       = module.user_pool.vite_user_pool_client_id
  description = "User pool client ID"
}

output "vite_user_pool_id" {
  value       = module.user_pool.vite_user_pool_id
  description = "Vite user pool ID"
}

output "vite_user_pool_client_id" {
  value       = module.user_pool.vite_user_pool_client_id
  description = "Vite user pool client ID"
}
