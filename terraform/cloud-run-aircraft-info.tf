resource "google_service_account" "aircraft_info_sa" {
  account_id   = "aircraft-info"
  display_name = "A service account for the Aircraft-info service"
}

resource "google_project_iam_binding" "aircraft_info_binding_log" {
  project = var.project_id
  role    = "roles/logging.logWriter"

  members = [
    "serviceAccount:${google_service_account.aircraft_info_sa.email}",
  ]
}

# Aircraft Info service
resource "google_cloud_run_service" "aircraft_info" {
  depends_on = [
    google_service_account.aircraft_info_sa,
    google_project_iam_binding.aircraft_info_binding_log,
    google_project_service.cloud_run,
  ]

  name     = "aircraft-info"
  location = var.region

  template {
    spec {
      containers {
        image = "${var.docker_repo_region}-docker.pkg.dev/${var.project_id}/${var.docker_repo_name}/aircraft_info:${var.aircraft_info_tag}"
        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = var.project_id
        }
      }
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "15"
        "autoscaling.knative.dev/minScale" = "0"
      }
    }
  }

  # direct all traffic toward latest revision
  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_policy" "noauth_aircraft_info" {
  location = google_cloud_run_service.aircraft_info.location
  project  = google_cloud_run_service.aircraft_info.project
  service  = google_cloud_run_service.aircraft_info.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
