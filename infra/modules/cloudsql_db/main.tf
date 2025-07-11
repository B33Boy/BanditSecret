# modules/cloudsql_db/main.tf

module "mysql_instance" {
  source  = "GoogleCloudPlatform/sql-db/google//modules/mysql"
  version = "~> 26.1"

  project_id       = var.project_id
  name             = var.cloudsql_instance_name
  database_version = var.database_version
  region           = var.region
  tier             = var.instance_tier

  ip_configuration = {
    ipv4_enabled                                  = true
    private_network                               = var.network_self_link # Input from network module's output
    allocated_ip_range                            = var.cloudsql_private_ip_range_name
    require_ssl                                   = false # Set to true if all clients enforce SSL/TLS
    enable_private_path_for_google_cloud_services = true
  }

  db_name = var.cloudsql_database_name

  enable_default_user = true
  additional_users = [
    {
      name            = var.cloudsql_username
      host            = "%" # Allows connection from any host (adjust for stricter security)
      random_password = true
      type            = "PRIMARY"
      password        = null
    }
  ]

  # Highly recommended to set it to true for prod
  deletion_protection_enabled = false

  depends_on = [
    var.api_dependencies,
    google_service_networking_connection.private_vpc_connection,
    # No direct dependency on module.vpc_network here because network_self_link is an input,
    # implying the network is already managed/created.
  ]
}

# Required for Private IP (Service Networking API and Connection)
# These resources are external to the Cloud SQL module but required for its private IP config.

# Allocate a private IP range for Google services in your VPC
resource "google_compute_global_address" "private_ip_alloc" {
  project       = var.project_id
  name          = var.cloudsql_private_ip_range_name
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 20                    # A /20 range is commonly used
  network       = var.network_self_link # Uses the network_self_link input
  depends_on    = [var.api_dependencies]
}

# Establish the private connection between your VPC and Google's services VPC
resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = var.network_self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
  depends_on = [
    var.api_dependencies,
    google_compute_global_address.private_ip_alloc
  ]
}

# resource "google_project_service" "servicenetworking_api" {
#   project            = var.project_id
#   service            = "servicenetworking.googleapis.com"
#   disable_on_destroy = false
# }
