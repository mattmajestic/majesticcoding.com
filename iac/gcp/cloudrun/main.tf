terraform {
  required_providers {
    google = { source = "hashicorp/google", version = "~> 5.0" }
  }
}

# --- set your project/region here ---
variable "project_id" { default = "YOUR_GCP_PROJECT_ID" }
variable "region"     { default = "us-central1" }

provider "google" {
  project = var.project_id
  region  = var.region
}

# Enable Cloud Run API
resource "google_project_service" "run" {
  project = var.project_id
  service = "run.googleapis.com"
}

# Cloud Run service (v2)
resource "google_cloud_run_v2_service" "app" {
  name     = "majesticcodingcom"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "docker.io/mattmajestic/majesticcodingcom:latest"
      # Ensure your container listens on $PORT (Cloud Run sets it).
      # env { name = "PORT", value = "8080" } # uncomment if your image needs it
    }
  }

  depends_on = [google_project_service.run]
}

# Make it publicly invokable
resource "google_cloud_run_v2_service_iam_member" "public" {
  project  = var.project_id
  location = var.region
  service  = google_cloud_run_v2_service.app.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

output "url" {
  value = google_cloud_run_v2_service.app.uri
}
