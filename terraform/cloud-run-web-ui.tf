# Web UI service
resource "google_cloud_run_service" "web_ui" {
  name     = "web-ui"
  location = var.region

  template {
    spec {
      containers {
<<<<<<< HEAD
        image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.docker_repo_name}/web_ui:latest"
=======
        image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/web_ui:latest"
>>>>>>> feature/dev-gateway
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

resource "google_cloud_run_service_iam_policy" "noauth_web_ui" {
  location = google_cloud_run_service.web_ui.location
  project  = google_cloud_run_service.web_ui.project
  service  = google_cloud_run_service.web_ui.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
