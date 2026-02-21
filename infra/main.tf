terraform {
  required_version = ">= 1.7.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.40"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_firestore_database" "default" {
  project                           = var.project_id
  name                              = "(default)"
  location_id                       = var.region
  type                              = "FIRESTORE_NATIVE"
  concurrency_mode                  = "OPTIMISTIC"
  app_engine_integration_mode       = "DISABLED"
  delete_protection_state           = "DELETE_PROTECTION_ENABLED"
  point_in_time_recovery_enablement = "POINT_IN_TIME_RECOVERY_ENABLED"
}

resource "google_storage_bucket" "menu_images" {
  name                        = "${var.project_id}-menu-images"
  location                    = var.region
  uniform_bucket_level_access = true
  force_destroy               = false
}

resource "google_service_account" "backend" {
  account_id   = "gourmet-guide-backend"
  display_name = "GourmetGuide Backend"
}

resource "google_cloud_run_v2_service" "backend" {
  name     = "gourmet-guide-backend"
  location = var.region

  template {
    service_account = google_service_account.backend.email
    scaling {
      min_instance_count = 0
      max_instance_count = 3
    }
    containers {
      image = var.backend_image
      resources {
        limits = {
          memory = "512Mi"
        }
      }
      env {
        name  = "GOOGLE_CLOUD_PROJECT"
        value = var.project_id
      }
      env {
        name  = "GEMINI_MODEL"
        value = var.gemini_model
      }
      env {
        name  = "GOOGLE_CLOUD_LOCATION"
        value = var.region
      }
      env {
        name  = "GOOGLE_GENAI_USE_VERTEXAI"
        value = "true"
      }
      env {
        name  = "MENU_IMAGE_BUCKET"
        value = google_storage_bucket.menu_images.name
      }
    }
    max_instance_request_concurrency = 80
  }

  depends_on = [google_firestore_database.default]
}
