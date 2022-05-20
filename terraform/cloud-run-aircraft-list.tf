
resource "google_service_account" "aircraft_list_sa" {
  account_id   = "aircraft-list"
  display_name = "A service account for the Aircraft List service"
}

# bind the service account to the necessary roles
resource "google_project_iam_binding" "aircraft_list_binding" {
  project = var.project_id
  role    = "roles/datastore.viewer"

  members = [
    "serviceAccount:${google_service_account.aircraft_list_sa.email}",
  ]
}

resource "google_service_account_key" "aircraft_list_key" {
  service_account_id = google_service_account.aircraft_list_sa.name
  public_key_type    = "TYPE_X509_PEM_FILE"
}

# Aircraft List service
resource "google_cloud_run_service" "aircraft_list" {
  depends_on = [
    google_service_account.aircraft_list_sa,
    google_project_iam_binding.aircraft_list_binding,
    google_service_account_key.aircraft_list_key,
  ]
  name     = "aircraft-list"
  location = var.region

  template {
    spec {
      # service_account_name = google_service_account.aircraft_list_sa.email
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/aircraft_list:latest"
        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = var.project_id
        }
        env {
          name  = "AUTHENTICATION_METHOD"
          value = "JSON"
        }
        env {
          name  = "GOOGLE_APPLICATION_CREDENTIALS"
          value = " ${base64decode(google_service_account_key.aircraft_list_key.private_key)} "
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
