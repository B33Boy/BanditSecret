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
  project     = var.gcp_project
  region      = var.gcp_region
  credentials = file(var.gcp_svc_key)
}
