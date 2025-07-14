# environments/dev/main.tf

# Enable the required APIs centrally
resource "google_project_service" "essential_apis" {
  for_each = toset([
    "compute.googleapis.com",
    "servicenetworking.googleapis.com",
    "vpcaccess.googleapis.com",
    "sqladmin.googleapis.com",
    "run.googleapis.com",
    "storage.googleapis.com",
    "eventarc.googleapis.com"
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

# --- GCS Buckets Module ---
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

# --- GCS Folder Objects ---
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


# --- Service Account Module ---
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


# --- Cloud Function: VTT To JSON --- 
module "vtt_to_json_converter" {
  source                  = "../../modules/cloud_function_gcs_trigger"
  name                    = "${var.project_id}-vtt-to-json"
  project_id              = var.project_id
  region                  = var.region
  runtime                 = "python312"
  entry_point             = "vtt_to_json_converter"
  source_code_bucket_name = module.gcs_buckets.buckets_map["function-sources"].name # Shared source bucket
  source_code_path        = "../../../cloud_functions/vtt_to_json_converter.zip"    # Local path to zipped source
  memory                  = "256MiB"
  timeout_seconds         = 60

  trigger_region    = var.region
  event_type        = "google.cloud.storage.object.v1.finalized"
  event_resource_id = module.gcs_buckets.buckets_map["captions"].id
  event_attribute_filters = {
    "name" : "raw_vtt/"
  }
}

# IAM bindings for vtt_to_json_converter
resource "google_storage_bucket_iam_member" "vtt_converter_gcs_reader" {
  bucket = module.gcs_buckets.buckets_map["captions"].name
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${module.vtt_to_json_converter.service_account_email}"
}

resource "google_storage_bucket_iam_member" "vtt_converter_gcs_creator" {
  bucket = module.gcs_buckets.buckets_map["captions"].name
  role   = "roles/storage.objectCreator"
  member = "serviceAccount:${module.vtt_to_json_converter.service_account_email}"
}

