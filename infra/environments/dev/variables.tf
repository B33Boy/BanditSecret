variable "project_id" {
  description = "The GCP project ID for the dev environment."
  type        = string
}

variable "region" {
  description = "The GCP region to deploy resources in."
  type        = string
  default     = "northamerica-northeast2"
}

variable "gcp_svc_key" {
  description = "Path to the GCP service account key JSON file used for Terraform authentication."
  type        = string
  sensitive   = true
}

# Network Variables
variable "vpc_network_name" {
  description = "The name of the VPC network."
  type        = string
  default     = "dev-vpc-network"
}

variable "vpc_subnet_name" {
  description = "The name of the primary subnet."
  type        = string
  default     = "primary"
}

variable "vpc_subnet_ip_cidr_range" {
  description = "The IP CIDR range for the primary subnet (e.g., '10.10.0.0/20')."
  type        = string
  default     = "10.10.0.0/28"
}

variable "vpc_connector_name" {
  description = "The name of the Serverless VPC Access connector."
  type        = string
  default     = "vpc-connector"
}

# Cloud SQL Variables
variable "cloudsql_instance_name" {
  description = "The Cloud SQL instance name."
  type        = string
  default     = "dev-app-cloudsql"
}

variable "database_version" {
  description = "The database version for Cloud SQL (e.g., MYSQL_8_0)."
  type        = string
  default     = "MYSQL_8_0"
}

variable "instance_tier" {
  description = "The machine tier for the Cloud SQL instance."
  type        = string
  default     = "db-f1-micro"
}

variable "cloudsql_private_ip_range_name" {
  description = "The name for the private IP address allocation."
  type        = string
  default     = "sql-private-range"
}

variable "cloudsql_database_name" {
  description = "The default database name to create in Cloud SQL."
  type        = string
  default     = "app_database"
}

variable "cloudsql_username" {
  description = "The username for the default Cloud SQL user."
  type        = string
  default     = "app_user"
}
