# Enable Cloud SQL Admin API
resource "google_project_service" "sqladmin_api" {
  project            = var.gcp_project
  service            = "sqladmin.googleapis.com"
  disable_on_destroy = false
}

# Enable Compute Engine API (needed for VPC, connectors)
resource "google_project_service" "compute_api" {
  project            = var.gcp_project
  service            = "compute.googleapis.com"
  disable_on_destroy = false
}

# Enable Service Networking API (for Private IP Cloud SQL)
resource "google_project_service" "servicenetworking_api" {
  project            = var.gcp_project
  service            = "servicenetworking.googleapis.com"
  disable_on_destroy = false
}

# Enable Serverless VPC Access API
resource "google_project_service" "vpcaccess_api" {
  project            = var.gcp_project
  service            = "vpcaccess.googleapis.com"
  disable_on_destroy = false
}

# Enable Cloud Run API
resource "google_project_service" "cloudrun_api" {
  project            = var.gcp_project
  service            = "run.googleapis.com"
  disable_on_destroy = false
}

# Enable Secret Manager API
resource "google_project_service" "secretmanager_api" {
  project            = var.gcp_project
  service            = "secretmanager.googleapis.com"
  disable_on_destroy = false
}
