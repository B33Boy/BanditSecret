# environments/dev/main.tf

# Enable the required APIs centrally
resource "google_project_service" "essential_apis" {
  for_each = toset([
    "compute.googleapis.com",
    "servicenetworking.googleapis.com",
    "vpcaccess.googleapis.com",
    "sqladmin.googleapis.com",
    "run.googleapis.com",
    "storage.googleapis.com"
  ])
  project            = var.project_id
  service            = each.key
  disable_on_destroy = false
}


# --- Network Module ---
module "network" {
  source = "../../modules/network"

  project_id               = var.project_id
  region                   = var.region
  vpc_network_name         = var.vpc_network_name
  vpc_subnet_name          = var.vpc_subnet_name
  vpc_subnet_ip_cidr_range = var.vpc_subnet_ip_cidr_range
  vpc_connector_name       = var.vpc_connector_name
}

# --- Cloud SQL Module ---
module "cloudsql_db" {
  source = "../../modules/cloudsql_db"

  project_id                     = var.project_id
  region                         = var.region
  cloudsql_instance_name         = var.cloudsql_instance_name
  database_version               = var.database_version
  instance_tier                  = var.instance_tier
  cloudsql_private_ip_range_name = var.cloudsql_private_ip_range_name
  cloudsql_database_name         = var.cloudsql_database_name
  cloudsql_username              = var.cloudsql_username
  network_self_link              = module.network.network_self_link
}


module "gcs_buckets" {
  source     = "terraform-google-modules/cloud-storage/google"
  version    = "~> 11.0"
  project_id = var.project_id
  location   = var.region

  names  = ["captions", "function-sources"]
  prefix = var.project_id

  force_destroy = {
    captions         = true
    function-sources = true
  }

  lifecycle_rules = [
    {
      bucket    = "captions"
      action    = { type = "Delete" }
      condition = { age = 30 }
    }
  ]
}


# Storage buckets don't have folder hierarchy internally, everythin is flat
# To simulate folder behaviour
resource "google_storage_bucket_object" "vtt_folder" {
  name    = "raw_vtt/"
  content = " "
  bucket  = module.gcs_buckets.buckets_map["captions"].name
}

resource "google_storage_bucket_object" "json_folder" {
  name    = "converted_json/"
  content = " "
  bucket  = module.gcs_buckets.buckets_map["captions"].name
}


module "service_account" {
  source     = "terraform-google-modules/service-accounts/google//modules/simple-sa"
  version    = "~> 4.0"
  project_id = var.project_id

  name         = "ytdlp-svc"
  display_name = "Service Account for yt-dlp uploads"
  project_roles = [
    "roles/viewer",
    "roles/storage.objectAdmin",
  ]
}


