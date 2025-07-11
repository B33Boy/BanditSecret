# modules/network/variables.tf

variable "project_id" {
  description = "The GCP project ID."
  type        = string
}

variable "region" {
  description = "The GCP region for the network resources."
  type        = string
}

variable "api_dependencies" {
  description = "A list of resources this module depends on to ensure APIs are enabled."
  type        = any
  default     = []
}

variable "vpc_network_name" {
  description = "The name of the VPC network."
  type        = string
}

variable "vpc_subnet_name" {
  description = "The name of the subnet within the VPC."
  type        = string
}

variable "vpc_subnet_ip_cidr_range" {
  description = "The IP CIDR range for the VPC subnet."
  type        = string
}

variable "vpc_connector_name" {
  description = "The name for the Serverless VPC Access Connector."
  type        = string
}
