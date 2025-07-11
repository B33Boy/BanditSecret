resource "google_service_account" "default" {
  project      = var.project_id
  account_id   = var.service_account_id
  display_name = var.display_name
}

output "email" {
  description = "The email of the created service account."
  value       = google_service_account.default.email
}

output "id" {
  description = "The ID of the created service account."
  value       = google_service_account.default.id
}
