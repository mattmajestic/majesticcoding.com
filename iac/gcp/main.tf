provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_cloud_run_service" "majesticcoding" {
  name     = "majesticcoding"
  location = var.region

  template {
    spec {
      containers {
        image = "docker.io/mattmajestic/majesticcoding:latest"
        # Add env vars if needed
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "invoker" {
  service  = google_cloud_run_service.majesticcoding.name
  location = google_cloud_run_service.majesticcoding.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

variable "project_id" {}
variable "region" { default = "us-central1" }