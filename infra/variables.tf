# Core
variable "gcp_project" {
  type = string
}

variable "gcp_region" {
  type    = string
  default = "northamerica-northeast2"
}

variable "gcp_svc_key" {
  type = string
}

# Cloud sql
variable "cloud_sql_instance_name" {
  type    = string
  default = "cloudsql-instance"
}

variable "cloud_sql_database_name" {
  type    = string
  default = "cloudsql-db"
}

variable "cloud_sql_user_name" {
  type    = string
  default = "cloudsql-user"
}

# Cloud run
variable "cloud_run_service_name" {
  type    = string
  default = "cloudrun-svc"
}

variable "cloud_run_image" {
  type    = string
  default = "cloudrun-image"
}

# Network
variable "vpc_network_name" {
  type    = string
  default = "vpc-net"
}

variable "vpc_subnetwork_name" {
  type    = string
  default = "vpc-subnet"
}

variable "vpc_subnetwork_cidr" {
  type    = string
  default = "10.0.0.0/20"
}

variable "vpc_connector_name" {
  type    = string
  default = "vpc-connector"
}

variable "vpc_connector_ip_range" {
  type    = string
  default = "10.0.16.0/28"
}
