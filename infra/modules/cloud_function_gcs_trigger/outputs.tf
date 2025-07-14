# modules/cloud_function_v2_gcs_trigger/outputs.tf

output "function_name" {
  description = "The name of the deployed Cloud Function."
  value       = google_cloudfunctions2_function.main.name
}

output "function_uri" {
  description = "The HTTP URI for the deployed Cloud Function (if applicable)."
  value       = google_cloudfunctions2_function.main.service_config.0.uri
}

output "service_account_email" {
  description = "The email of the service account created for the Cloud Function."
  value       = google_service_account.function_sa.email
}

output "service_account_id" {
  description = "The ID of the service account created for the Cloud Function."
  value       = google_service_account.function_sa.id
}
