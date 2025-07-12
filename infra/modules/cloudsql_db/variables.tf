# modules/cloudsql_db/variables.tf

variable "project_id" {
  description = "The GCP project ID."
  type        = string
}

variable "region" {
  description = "The region to deploy Cloud SQL resources."
  type        = string
}

variable "cloudsql_instance_name" {
  description = "Name for the Cloud SQL instance."
  type        = string
}

variable "database_version" {
  description = "Database engine version (e.g., MYSQL_8_0)."
  type        = string
}

variable "instance_tier" {
  description = "The tier (machine type) for Cloud SQL (e.g., db-f1-micro)."
  type        = string
}

variable "cloudsql_private_ip_range_name" {
  description = "Name of the IP range used for private service access."
  type        = string
}

variable "cloudsql_database_name" {
  description = "The name of the default database."
  type        = string
}

variable "cloudsql_username" {
  description = "Username for the default user."
  type        = string
}

variable "network_self_link" {
  description = "The self_link of the VPC network."
  type        = string
}
