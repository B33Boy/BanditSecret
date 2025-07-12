# modules/cloudsql_db/main.tf

resource "google_compute_global_address" "private_ip_alloc" {
  project       = var.project_id
  name          = var.cloudsql_private_ip_range_name
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 20
  network       = var.network_self_link
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = var.network_self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]

  depends_on = [google_compute_global_address.private_ip_alloc]
}

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
    private_network                               = var.network_self_link
    allocated_ip_range                            = var.cloudsql_private_ip_range_name
    require_ssl                                   = false
    enable_private_path_for_google_cloud_services = true
  }

  db_name = var.cloudsql_database_name

  enable_default_user = true
  additional_users = [
    {
      name            = var.cloudsql_username
      host            = "%"
      random_password = true
      type            = ""
      password        = null
    }
  ]

  deletion_protection_enabled = false

  depends_on = [google_service_networking_connection.private_vpc_connection]
}
