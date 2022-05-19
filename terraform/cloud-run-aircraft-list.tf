
resource "google_service_account" "aircraft_list_sa" {
  account_id   = "aircraft-list"
  display_name = "A service account for the Aircraft List service"
}

# bind the service account to the necessary roles
resource "google_project_iam_binding" "aircraft_list_binding" {
  project = var.project_id
  role    = "roles/datastore.user"

  members = [
    "serviceAccount:${google_service_account.aircraft_list_sa.email}",
  ]
}

# Aircraft List service
resource "google_cloud_run_service" "aircraft_list" {
  name     = "aircraft-list"
  location = var.region

  template {
    spec {
      service_account_name = google_service_account.aircraft_list_sa.email
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/aircraft_list:latest"
        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = var.project_id
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "2"
        "autoscaling.knative.dev/minScale" = "1"
      }
    }
  }

  # direct all traffic toward latest revision
  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_policy" "noauth_aircraft_list" {
  location = google_cloud_run_service.aircraft_list.location
  project  = google_cloud_run_service.aircraft_list.project
  service  = google_cloud_run_service.aircraft_list.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
