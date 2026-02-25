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

locals {
  name_prefix      = "gourmet-guide-${var.environment}"
  image_bucket     = "${var.project_id}-${var.environment}-menu-images"
  frontend_bucket  = "${var.project_id}-${var.environment}-frontend"
  run_service_name = "${local.name_prefix}-backend"
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
  name                        = local.image_bucket
  location                    = var.region
  uniform_bucket_level_access = true
  force_destroy               = false
}

resource "google_storage_bucket" "frontend" {
  name                        = local.frontend_bucket
  location                    = var.region
  uniform_bucket_level_access = true
  force_destroy               = false

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
}

resource "google_storage_bucket_iam_member" "frontend_public" {
  count = var.frontend_public_access ? 1 : 0

  bucket = google_storage_bucket.frontend.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

resource "google_service_account" "backend_runtime" {
  account_id   = "${replace(local.name_prefix, "-", "")}-runtime"
  display_name = "GourmetGuide ${var.environment} Backend Runtime"
}

resource "google_service_account" "backend_deployer" {
  account_id   = "${replace(local.name_prefix, "-", "")}-deployer"
  display_name = "GourmetGuide ${var.environment} Backend Deployer"
}

resource "google_project_iam_member" "runtime_firestore" {
  project = var.project_id
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.backend_runtime.email}"
}

resource "google_storage_bucket_iam_member" "runtime_storage" {
  bucket = google_storage_bucket.menu_images.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.backend_runtime.email}"
}

resource "google_project_iam_member" "deployer_run_admin" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${google_service_account.backend_deployer.email}"
}

resource "google_project_iam_member" "deployer_sa_user" {
  project = var.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.backend_deployer.email}"
}

resource "google_cloud_run_v2_service" "backend" {
  name     = local.run_service_name
  location = var.region
  ingress  = var.cloud_run_ingress

  template {
    service_account = google_service_account.backend_runtime.email
    scaling {
      min_instance_count = 0
      max_instance_count = 3
    }

    dynamic "vpc_access" {
      for_each = var.vpc_connector != "" ? [1] : []
      content {
        connector = var.vpc_connector
        egress    = "PRIVATE_RANGES_ONLY"
      }
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

  depends_on = [
    google_firestore_database.default,
    google_project_iam_member.runtime_firestore,
    google_storage_bucket_iam_member.runtime_storage,
  ]
}

resource "google_cloud_run_service_iam_member" "public_invoker" {
  count    = var.allow_unauthenticated ? 1 : 0
  location = google_cloud_run_v2_service.backend.location
  project  = var.project_id
  service  = google_cloud_run_v2_service.backend.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
