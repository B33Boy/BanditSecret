variable "project_id" {
  description = "The GCP project ID."
  type        = string
}

variable "service_account_id" {
  description = "The ID of the service account to create."
  type        = string
}

variable "display_name" {
  description = "The display name for the service account."
  type        = string
}
