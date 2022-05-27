# resource "google_service_account" "aircraft_list_sa" {
#   account_id   = "aircraft-list"
#   display_name = "A service account for the Aircraft List service"
# }

# # bind the service account to the necessary roles
# resource "google_project_iam_binding" "aircraft_list_binding" {
#   project = var.project_id
#   role    = "roles/datastore.user"

#   members = [
#     "serviceAccount:${google_service_account.aircraft_list_sa.email}",
#   ]
# }

locals {
  # aircraft_list_sa_name = google_service_account.aircraft_list_sa.name
  # aircraft_list_sa_email = google_service_account.aircraft_list_sa.email
  aircraft_list_sa_email = "aircraft-list@${var.project_id}.iam.gserviceaccount.com"
  aircraft_list_sa_name  = "projects/${var.project_id}/serviceAccounts/${local.aircraft_list_sa_email}"
}

# Aircraft List service
resource "google_cloud_run_service" "aircraft_list" {
  depends_on = [
    # google_service_account.aircraft_list_sa,

  ]
  name     = "aircraft-list"
  location = var.region

  template {
    spec {
      service_account_name  = local.aircraft_list_sa_email
      container_concurrency = 1000
      containers {
        image = "${var.docker_repo_region}-docker.pkg.dev/${var.project_id}/${var.docker_repo_name}/aircraft_list:${var.aircraft_list_tag}"
        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = var.project_id
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "10"
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
