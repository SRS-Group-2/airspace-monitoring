# resource "google_service_account" "airspace_monthly_history_sa" {
#   account_id   = "airspace-monthly-history"
#   display_name = "A service account for the Airspace Monthly History service"
# }

# # bind the service account to the necessary roles
# resource "google_project_iam_binding" "airspace_monthly_history_binding" {
#   project = var.project_id
#   role    = "roles/datastore.user"

#   members = [
#     "serviceAccount:${google_service_account.airspace_monthly_history_sa.email}",
#   ]
# }

locals {
  # airspace_monthly_history_sa_name = google_service_account.airspace_monthly_history_sa.name
  # airspace_monthly_history_sa_email = google_service_account.airspace_monthly_history_sa.email
  airspace_monthly_history_sa_email = "airspace-monthly-history@${var.project_id}.iam.gserviceaccount.com"
  airspace_monthly_history_sa_name  = "projects/${var.project_id}/serviceAccounts/${local.airspace_monthly_history_sa_email}"
}

# resource "google_service_account_key" "airspace_monthly_history_key" {
#   service_account_id = local.airspace_monthly_history_sa_name
#   public_key_type    = "TYPE_X509_PEM_FILE"
# }

# Airspace Monthly History service
resource "google_cloud_run_service" "airspace_monthly_history" {
  depends_on = [
    # google_service_account_key.airspace_monthly_history_key,
  ]
  name     = "airspace-monthly-history"
  location = var.region

  template {
    spec {
      service_account_name = local.airspace_monthly_history_sa_email
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.docker_repo_name}/airspace_monthly_history:latest"
        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = var.project_id
        }
        env {
          name  = "AUTHENTICATION_METHOD"
          value = "ADC"
        }
        # env {
        #   name  = "GOOGLE_APPLICATION_CREDENTIALS"
        #   value = " ${base64decode(google_service_account_key.airspace_monthly_history_key.private_key)} "
        # }
        env {
          name  = "GIN_MODE"
          value = "release"
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "2"
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

resource "google_cloud_run_service_iam_policy" "noauth_airspace_monthly_history" {
  location = google_cloud_run_service.airspace_monthly_history.location
  project  = google_cloud_run_service.airspace_monthly_history.project
  service  = google_cloud_run_service.airspace_monthly_history.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
