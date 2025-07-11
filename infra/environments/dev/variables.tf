# environments/dev/variables.tf
variable "project_id" {
  description = "The GCP project ID for the dev environment."
  type        = string
}

variable "region" {
  description = "The GCP region to deploy resources in."
  type        = string
  default     = "northamerica-northeast2"
}

# Service Account Key Path (Sensitive!)
variable "gcp_svc_key" {
  description = "Path to the GCP service account key JSON file used for Terraform authentication."
  type        = string
  sensitive   = true
}

# Network Variables (using prefixes/suffixes for modular naming)
variable "vpc_network_name_prefix" {
  description = "Prefix for the VPC network name. Will be combined with project_id."
  type        = string
  default     = "dev-vpc-network"
}

variable "vpc_subnet_name_suffix" {
  description = "Suffix for the primary subnet name. Will be combined with network name."
  type        = string
  default     = "primary"
}

variable "vpc_subnet_ip_cidr_range" {
  description = "The IP CIDR range for the primary subnet (e.g., '10.10.0.0/20')."
  type        = string
  default     = "10.10.0.0/20"
}

variable "vpc_connector_name_suffix" {
  description = "Suffix for the Serverless VPC Access connector name."
  type        = string
  default     = "connector"
}

# Cloud SQL Variables (using prefixes/suffixes for modular naming)
variable "cloudsql_instance_name_prefix" {
  description = "Prefix for the Cloud SQL instance name. Will be combined with project_id."
  type        = string
  default     = "dev-app-cloudsql"
}

variable "database_version" {
  description = "The database version for Cloud SQL (e.g., MYSQL_8_0)."
  type        = string
  default     = "MYSQL_8_0"
}

variable "instance_tier" {
  description = "The machine tier for the Cloud SQL instance (e.g., db-f1-micro, db-g1-small)."
  type        = string
  default     = "db-f1-micro"
}

variable "cloudsql_private_ip_range_name_suffix" {
  description = "Suffix for the global address resource name for Cloud SQL private service access."
  type        = string
  default     = "sql-private-range"
}

variable "cloudsql_database_name" {
  description = "The name of the default database in Cloud SQL."
  type        = string
  default     = "app_database"
}

variable "cloudsql_username" {
  description = "The username for the default Cloud SQL user."
  type        = string
  default     = "app_user"
}
