# resource "google_service_account" "airspace_daily_history_sa" {
#   account_id   = "airspace-daily-history"
#   display_name = "A service account for the Airspace Daily History service"
# }

# # bind the service account to the necessary roles
# resource "google_project_iam_binding" "airspace_daily_history_binding" {
#   project = var.project_id
#   role    = "roles/datastore.viewer"

#   members = [
#     "serviceAccount:${google_service_account.airspace_daily_history_sa.email}",
#   ]
# }

locals {
  # airspace_daily_history_sa_name = google_service_account.airspace_daily_history_sa.name
  # airspace_daily_history_sa_email = google_service_account.airspace_daily_history_sa.email
  airspace_daily_history_sa_email = "airspace-daily-history@${var.project_id}.iam.gserviceaccount.com"
  airspace_daily_history_sa_name  = "projects/${var.project_id}/serviceAccounts/${local.airspace_daily_history_sa_email}"
}

# Aircraft Daily History service
resource "google_cloud_run_service" "airspace_daily_history" {
  depends_on = [
    # google_service_account.airspace_daily_history_sa,
  ]
  name     = "airspace-daily-history"
  location = var.region

  template {
    spec {
      service_account_name  = local.airspace_daily_history_sa_email
      container_concurrency = 100
      containers {
        image = "${var.docker_repo_region}-docker.pkg.dev/${var.docker_repo_project_id}/${var.docker_repo_name}/airspace_daily_history:${var.airspace_daily_history_tag}"
        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = var.project_id
        }
        env {
          name  = "GIN_MODE"
          value = "release"
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

resource "google_cloud_run_service_iam_policy" "noauth_airspace_daily_history" {
  location = google_cloud_run_service.airspace_daily_history.location
  project  = google_cloud_run_service.airspace_daily_history.project
  service  = google_cloud_run_service.airspace_daily_history.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
