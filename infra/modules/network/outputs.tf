# modules/network/outputs.tf

output "network_self_link" {
  description = "Self link of the created VPC network."
  value       = module.vpc_network.network_self_link
}

output "vpc_connector_id" {
  description = "ID of the Serverless VPC Access connector."
  value       = module.serverless_connector.connector_ids
}
