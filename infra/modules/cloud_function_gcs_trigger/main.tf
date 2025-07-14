# modules/cloud_function_v2_gcs_trigger/main.tf


resource "google_service_account" "function_sa" {
  account_id   = "${var.name}-svc"
  display_name = "Service Account for Cloud Functions: ${var.name}"
  project      = var.project_id
}

resource "google_storage_bucket_object" "function_source" {
  name   = "${var.name}-source-${timestamp()}.zip"
  bucket = var.source_code_bucket_name
  source = var.source_code_path
}

/*
  Cloud Functions 2nd Gen uses Eventarc, 
  and Eventarc relies on Cloud Storage events being published via Pub/Sub behind the scenes.
  So we need to give the internal GCS service account the role/pubsub.publisher 
*/
data "google_project" "current" {
  project_id = var.project_id
}

resource "google_project_iam_member" "gcs_pubsub_publisher" {
  project = var.project_id
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:service-${data.google_project.current.number}@gs-project-accounts.iam.gserviceaccount.com"
}


resource "google_cloudfunctions2_function" "main" {
  name     = var.name
  location = var.region
  project  = var.project_id

  build_config {
    runtime     = var.runtime
    entry_point = var.entry_point
    source {
      storage_source {
        bucket = google_storage_bucket_object.function_source.bucket
        object = google_storage_bucket_object.function_source.name
      }
    }
  }

  service_config {
    service_account_email = google_service_account.function_sa.email
    available_memory      = var.memory
    timeout_seconds       = var.timeout_seconds
    environment_variables = var.environment_variables
    max_instance_count    = var.max_instance_count
    min_instance_count    = var.min_instance_count
  }

  event_trigger {
    trigger_region = var.trigger_region
    event_type     = var.event_type
    event_filters {
      attribute = "bucket"
      value     = var.event_resource_id
    }

    dynamic "event_filters" {
      for_each = var.event_attribute_filters
      content {
        attribute = event_filters.key
        value     = event_filters.value
      }
    }
  }
}
