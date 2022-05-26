# cloud run
resource "google_project_service" "cloud_run" {
  project = var.project_id
  service = "run.googleapis.com"

  disable_dependent_services = true
}

# kubernetes
resource "google_project_service" "kubernetes" {
  project = var.project_id
  service = "container.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "compute" {
  project = var.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "autoscaling" {
  project = var.project_id
  service = "autoscaling.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "pubsub" {
  project = var.project_id
  service = "pubsub.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "firestore" {
  project = var.project_id
  service = "firestore.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "api_gateway" {
  project = var.project_id
  service = "apigateway.googleapis.com"

  disable_dependent_services = true
}
