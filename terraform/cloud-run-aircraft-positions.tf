# resource "google_service_account" "aircraft_positions_sa" {
#   account_id   = "aircraft-positions"
#   display_name = "A service account for the Aircraft Positions service"
# }

# # bind the service account to the necessary roles
# resource "google_project_iam_binding" "aircraft_positions_binding_pubsub" {
#   project = var.project_id
#   role    = "roles/pubsub.editor"

#   members = [
#     "serviceAccount:${google_service_account.aircraft_positions_sa.email}",
#   ]
# }

# resource "google_project_iam_binding" "aircraft_positions_binding_log" {
#   project = var.project_id
#   role    = "roles/logging.logWriter"

#   members = [
#     "serviceAccount:${google_service_account.aircraft_positions_sa.email}",
#   ]
# }

locals {
  # aircraft_positions_sa_name = google_service_account.aircraft_positions_sa.name
  # aircraft_positions_sa_email = google_service_account.aircraft_positions_sa.email
  aircraft_positions_sa_email = "aircraft-positions@${var.project_id}.iam.gserviceaccount.com"
  aircraft_positions_sa_name  = "projects/${var.project_id}/serviceAccounts/${local.aircraft_positions_sa_email}"
}

# Aircraft Daily History service
resource "google_cloud_run_service" "aircraft_positions" {
  # depends_on = [
  #   google_project_iam_binding.aircraft_positions_binding_log,
  #   google_project_iam_binding.aircraft_positions_binding_pubsub,
  # ]
  name     = "aircraft-positions"
  location = var.region

  template {
    spec {
      service_account_name  = local.aircraft_positions_sa_email
      container_concurrency = 100
      containers {
        image = "${var.docker_repo_region}-docker.pkg.dev/${var.docker_repo_project_id}/${var.docker_repo_name}/aircraft_positions:${var.aircraft_positions_tag}"
        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = var.project_id
        }
        env {
          name  = "GOOGLE_PUBSUB_AIRCRAFT_POSITIONS_TOPIC_ID"
          value = var.vectors_topic
        }
        env {
          name  = "GIN_MODE"
          value = "release"
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "10"
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

resource "google_cloud_run_service_iam_policy" "noauth_aircraft_positions" {
  location = google_cloud_run_service.aircraft_positions.location
  project  = google_cloud_run_service.aircraft_positions.project
  service  = google_cloud_run_service.aircraft_positions.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
