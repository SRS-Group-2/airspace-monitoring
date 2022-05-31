# Aircraft List service
resource "google_cloud_run_service" "websocket_endpoints" {
  depends_on = [
    google_cloud_run_service.aircraft_positions,
  ]
  name     = "websocket-endpoints"
  location = var.region

  template {
    spec {
      container_concurrency = 100
      containers {
        image = "${var.docker_repo_region}-docker.pkg.dev/${var.docker_repo_project_id}/${var.docker_repo_name}/websocket_endpoints:${var.websocket_endpoints_tag}"
        env {
          name  = "WEBSOCKET_ENDPOINT"
          value = google_cloud_run_service.aircraft_positions.status[0].url
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "5"
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

resource "google_cloud_run_service_iam_policy" "noauth_websocket_endpoints" {
  location = google_cloud_run_service.websocket_endpoints.location
  project  = google_cloud_run_service.websocket_endpoints.project
  service  = google_cloud_run_service.websocket_endpoints.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
