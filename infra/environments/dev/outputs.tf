# environments/dev/outputs.tf

# Network Outputs
output "vpc_network_self_link" {
  description = "The self_link of the VPC network."
  value       = module.network.network_self_link
}

output "vpc_connector_id" {
  description = "The ID of the Serverless VPC Access connector."
  value       = module.network.vpc_connector_id
}

# CloudSQL Outputs
output "cloudsql_connection_name" {
  description = "The connection name of the Cloud SQL instance."
  value       = module.cloudsql_db.instance_connection_name
}

output "cloudsql_private_ip" {
  description = "The private IP address of the Cloud SQL instance."
  value       = module.cloudsql_db.private_ip_address
}

output "cloudsql_user_password" {
  description = "The generated password for the Cloud SQL user."
  sensitive   = true
  value       = module.cloudsql_db.generated_user_password
}

# Cloud Storage Outputs
output "gcs_bucket_name" {
  value = google_storage_bucket.caption_bucket.name
}
