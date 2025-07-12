# modules/cloudsql_db/outputs.tf

output "instance_connection_name" {
  description = "Connection name of the Cloud SQL instance."
  value       = module.mysql_instance.instance_connection_name
}

output "private_ip_address" {
  description = "Private IP address of the Cloud SQL instance."
  value       = module.mysql_instance.private_ip_address
}

output "generated_user_password" {
  description = "Auto-generated password for the default user."
  sensitive   = true
  value       = module.mysql_instance.generated_user_password
}
