# environments/dev/main.tf

# Enable the required APIs centrally
resource "google_project_service" "essential_apis" {
  for_each = toset([
    "compute.googleapis.com",
    "servicenetworking.googleapis.com",
    "vpcaccess.googleapis.com",
    "sqladmin.googleapis.com",
    "run.googleapis.com"
  ])
  project = var.project_id
  service = each.key
}


# --- Network Module ---
module "network" {
  source = "../../modules/network"

  project_id               = var.project_id
  region                   = var.region
  api_dependencies         = google_project_service.essential_apis
  vpc_network_name         = var.vpc_network_name_prefix
  vpc_subnet_name          = var.vpc_subnet_name_suffix
  vpc_subnet_ip_cidr_range = var.vpc_subnet_ip_cidr_range
  vpc_connector_name       = var.vpc_connector_name_suffix
}

# --- Cloud SQL Module ---
module "cloudsql_db" {
  source = "../../modules/cloudsql_db"

  project_id                     = var.project_id
  region                         = var.region
  api_dependencies               = google_project_service.essential_apis
  cloudsql_instance_name         = var.cloudsql_instance_name_prefix
  database_version               = var.database_version
  instance_tier                  = var.instance_tier
  cloudsql_private_ip_range_name = var.cloudsql_private_ip_range_name_suffix
  cloudsql_database_name         = var.cloudsql_database_name
  cloudsql_username              = var.cloudsql_username
  network_self_link              = module.network.network_self_link
}

# Call a Cloud Run module here, referencing the VPC connector and Cloud SQL outputs
# For example:
/*
module "cloud_run_service" {
  source = "../../modules/cloud_run_app" # Example: Create a new module for your Cloud Run app
  project_id = var.project_id
  region     = var.region
  service_name = "my-basic-app" # Or use a variable for this

  # Define other Cloud Run specific inputs
  image = "gcr.io/${var.project_id}/your-app-image:latest" # Replace with your app image

  # Configure Cloud Run to use the VPC connector
  vpc_connector_id = module.network.vpc_connector_id # Output from network module

  # Pass Cloud SQL connection details to Cloud Run (e.g., via environment variables)
  # IMPORTANT: Manage sensitive data like passwords securely (e.g., Secret Manager)
  # For basic setup, you might pass them directly, but for production, avoid this.
  env_vars = {
    DB_HOST         = module.cloudsql_db.private_ip_address
    DB_NAME         = var.cloudsql_database_name
    DB_USER         = var.cloudsql_username
    DB_PASSWORD     = module.cloudsql_db.generated_user_password # Sensitive!
    INSTANCE_CONNECTION_NAME = module.cloudsql_db.instance_connection_name
  }
}
*/
