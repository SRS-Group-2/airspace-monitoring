# define Terraform general information
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.5.0"
    }
  }

  # where to save states
  backend "gcs" {
    bucket = "choir-tf-state"
    prefix = "terraform/state"
  }
}

# define cloud provider
provider "google" {
  project = var.project_id
  region  = var.region
  zone    = "${var.region}-c"
}

# noauth policy, necessary to call Cloud Run services as unauthenticated users
data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

# Aircraft Info service
resource "google_cloud_run_service" "aircraft_info" {
  name     = "aircraft_info"
  location = var.region

  template {
    spec {
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/aircraft_info:latest"
      }
    }
  }

  # direct all traffic toward latest revision
  traffic {
    percent         = 100
    latest_revision = true
  }

  metadata {
    annotations = {
      "autoscaling.knative.dev/maxScale" = "15"
      "autoscaling.knative.dev/minScale" = "0"
    }
  }
}

resource "google_cloud_run_service_iam_policy" "noauth_aircraft_info" {
  location = google_cloud_run_service.aircraft_info.location
  project  = google_cloud_run_service.aircraft_info.project
  service  = google_cloud_run_service.aircraft_info.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

# Aircraft List service
resource "google_cloud_run_service" "aircraft_list" {
  name     = "aircraft_list"
  location = var.region

  template {
    spec {
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/aircraft_list:latest"
      }
    }
  }

  # direct all traffic toward latest revision
  traffic {
    percent         = 100
    latest_revision = true
  }

  metadata {
    annotations = {
      "autoscaling.knative.dev/maxScale" = "2"
      "autoscaling.knative.dev/minScale" = "0"
    }
  }
}

resource "google_cloud_run_service_iam_policy" "noauth_aircraft_list" {
  location = google_cloud_run_service.aircraft_list.location
  project  = google_cloud_run_service.aircraft_list.project
  service  = google_cloud_run_service.aircraft_list.name

  policy_data = data.google_iam_policy.noauth.policy_data
}