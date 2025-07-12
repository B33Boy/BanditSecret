# environments/dev/providers.tf
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.42.0"
    }
  }

  backend "gcs" {
    bucket = "terraform-state-banditsecret"
    prefix = "terraform/state"
  }
}

provider "google" {
  project     = var.project_id
  region      = var.region
  credentials = file(var.gcp_svc_key)
}
