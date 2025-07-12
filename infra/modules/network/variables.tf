# modules/network/variables.tf

variable "project_id" {
  description = "GCP project ID."
  type        = string
}

variable "region" {
  description = "GCP region for VPC resources."
  type        = string
}

variable "vpc_network_name" {
  description = "Name of the VPC network."
  type        = string
}

variable "vpc_subnet_name" {
  description = "Name of the subnet in the VPC."
  type        = string
}

variable "vpc_subnet_ip_cidr_range" {
  description = "CIDR range for the subnet."
  type        = string
}

variable "vpc_connector_name" {
  description = "Name for the Serverless VPC Access connector."
  type        = string
}
