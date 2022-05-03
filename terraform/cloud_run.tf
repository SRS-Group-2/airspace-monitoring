
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
  name     = "aircraft-info"
  location = var.region

  template {
    spec {
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/aircraft_info:latest"
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

# Aircraft List service
resource "google_cloud_run_service" "aircraft_list" {
  name     = "aircraft-list"
  location = var.region

  template {
    spec {
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/aircraft_list:latest"
        env {
          name  = "GOOGLE_PUBSUB_AIRCRAFT_LIST_TOPIC_ID"
          value = var.aircraft_list_topic
        }
        env {
          name  = "GOOGLE_CREDENTIALS"
          value = var.google_list_credentials
        }
        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = var.project_id
        }
        env {
          name  = "GOOGLE_PUBSUB_AIRCRAFT_LIST_SUBSCRIBER_ID"
          value = var.aircraft_list_subscriber_id
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

resource "google_cloud_run_service_iam_policy" "noauth_aircraft_list" {
  location = google_cloud_run_service.aircraft_list.location
  project  = google_cloud_run_service.aircraft_list.project
  service  = google_cloud_run_service.aircraft_list.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
