resource "google_project_service" "artifact_registry" {
  project = var.project_id
  service = "artifactregistry.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "autoscaling" {
  project = var.project_id
  service = "autoscaling.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "compute" {
  project = var.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "container" {
  project = var.project_id
  service = "container.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "container_filesystem" {
  project = var.project_id
  service = "containerfilesystem.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "datastore" {
  project = var.project_id
  service = "datastore.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "firestore" {
  project = var.project_id
  service = "firestore.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "iam" {
  project = var.project_id
  service = "iam.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "iam_credentials" {
  project = var.project_id
  service = "iamcredentials.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "monitoring" {
  project = var.project_id
  service = "monitoring.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "pubsub" {
  project = var.project_id
  service = "pubsub.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "run" {
  project = var.project_id
  service = "run.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "storage_api" {
  project = var.project_id
  service = "storage-api.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "storage_component" {
  project = var.project_id
  service = "storage-component.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "storage" {
  project = var.project_id
  service = "storage.googleapis.com"

  disable_dependent_services = true
}
