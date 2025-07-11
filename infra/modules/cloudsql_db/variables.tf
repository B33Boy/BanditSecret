# modules/cloudsql_db/variables.tf

variable "project_id" {
  description = "The GCP project ID."
  type        = string
}

variable "api_dependencies" {
  description = "A list of resources this module depends on to ensure APIs are enabled."
  type        = any
  default     = []
}

variable "cloudsql_instance_name" {
  description = "Name for the Cloud SQL instance."
  type        = string
}

variable "database_version" {
  description = "The database version to use for the Cloud SQL instance (e.g., MYSQL_8_0)."
  type        = string
}

variable "region" {
  description = "The GCP region where the Cloud SQL instance will be created."
  type        = string
}

variable "instance_tier" {
  description = "The machine type (tier) for the Cloud SQL instance (e.g., db-f1-micro, db-g1-small)."
  type        = string
}

variable "cloudsql_private_ip_range_name" {
  description = "Name of the global address resource to be used for private service access (VPC Peering)."
  type        = string
}

variable "cloudsql_database_name" {
  description = "The name of the default database to create within the Cloud SQL instance."
  type        = string
}

variable "cloudsql_username" {
  description = "The username for the default Cloud SQL user."
  type        = string
}

variable "network_self_link" {
  description = "The self_link of the VPC network where Cloud SQL will be peered."
  type        = string
}
