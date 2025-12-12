provider "google" {
  project = var.project_id
  region  = var.region
}

variable "project_id" {
  description = "GCP Project ID"
}

variable "region" {
  description = "GCP Region"
  default     = "us-central1"
}

# Enable APIs
resource "google_project_service" "apis" {
  for_each = toset([
    "run.googleapis.com",
    "firestore.googleapis.com",
    "aiplatform.googleapis.com",
    "dlp.googleapis.com",
    "iam.googleapis.com",
  ])
  service = each.key
  disable_on_destroy = false
}

# Firestore Database
resource "google_firestore_database" "database" {
  name        = "(default)"
  location_id = var.region
  type        = "FIRESTORE_NATIVE"
  depends_on  = [google_project_service.apis]
}

# Cloud Run Service
resource "google_cloud_run_v2_service" "default" {
  name     = "prompt-defender"
  location = var.region
  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "gcr.io/${var.project_id}/prompt-defender:latest" // Placeholder image
      env {
        name  = "GOOGLE_PROJECT_ID"
        value = var.project_id
      }
      env {
        name  = "GOOGLE_LOCATION"
        value = var.region
      }
      env {
        name = "CACHE_COLLECTION_NAME"
        value = "cache"
      }
      env {
        name = "USERS_COLLECTION_NAME"
        value = "users"
      }
    }
  }
  depends_on = [google_project_service.apis]
}

# Service Account for Cloud Run (Default SA used above, but explicit is better)
# Granting roles
resource "google_project_iam_member" "firestore_user" {
  project = var.project_id
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_cloud_run_v2_service.default.template[0].service_account}"
}

resource "google_project_iam_member" "vertex_user" {
  project = var.project_id
  role    = "roles/aiplatform.user"
  member  = "serviceAccount:${google_cloud_run_v2_service.default.template[0].service_account}"
}

resource "google_project_iam_member" "dlp_user" {
  project = var.project_id
  role    = "roles/dlp.user"
  member  = "serviceAccount:${google_cloud_run_v2_service.default.template[0].service_account}"
}
