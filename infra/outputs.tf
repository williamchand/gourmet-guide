output "cloud_run_url" {
  description = "Backend Cloud Run URL"
  value       = google_cloud_run_v2_service.backend.uri
}

output "cloud_run_service_name" {
  description = "Cloud Run service name"
  value       = google_cloud_run_v2_service.backend.name
}

output "menu_image_bucket" {
  description = "Bucket for menu and dish images"
  value       = google_storage_bucket.menu_images.name
}

output "runtime_service_account_email" {
  description = "Runtime service account for backend access"
  value       = google_service_account.backend_runtime.email
}

output "deployer_service_account_email" {
  description = "Deployment service account email"
  value       = google_service_account.backend_deployer.email
}
