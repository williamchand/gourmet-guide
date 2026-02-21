output "cloud_run_url" {
  description = "Backend Cloud Run URL"
  value       = google_cloud_run_v2_service.backend.uri
}

output "menu_image_bucket" {
  description = "Bucket for menu and dish images"
  value       = google_storage_bucket.menu_images.name
}
