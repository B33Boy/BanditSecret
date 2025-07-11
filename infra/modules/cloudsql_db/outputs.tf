# modules/cloudsql_db/outputs.tf

output "instance_connection_name" {
  description = "The connection name of the master instance to be used in connection strings."
  value       = module.mysql_instance.instance_connection_name
}

output "private_ip_address" {
  description = "The private IPv4 address assigned for the master instance."
  value       = module.mysql_instance.private_ip_address
}

output "generated_user_password" {
  description = "The randomly generated password for the default user."
  sensitive   = true # Mark as sensitive to prevent plaintext logging
  value       = module.mysql_instance.generated_user_password
}
