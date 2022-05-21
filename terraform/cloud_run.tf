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
  depends_on = [google_cloud_run_service.aircraft_info]

  location = google_cloud_run_service.aircraft_info.location
  project  = google_cloud_run_service.aircraft_info.project
  service  = google_cloud_run_service.aircraft_info.name

  policy_data = data.google_iam_policy.noauth.policy_data
}


# # Aircraft List service
resource "google_cloud_run_service" "aircraft_list" {
  name     = "aircraft-list"
  location = var.region

  template {
    spec {
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/aircraft_list:latest"

        env {
          name  = "GOOGLE_APPLICATION_CREDENTIALS"
          value = var.google_list_credentials
        }
        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = var.project_id
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

# # Aircraft History service
resource "google_cloud_run_service" "airspace_daily_history" {
  name     = "airspace-daily-history"
  location = var.region

  template {
    spec {
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/airspace_daily_history:latest"
        env {
          name  = "GOOGLE_APPLICATION_CREDENTIALS"
          value = var.google_list_credentials
        }
        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = var.project_id
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

resource "google_cloud_run_service_iam_policy" "noauth_airspace_daily_history" {
  location    = google_cloud_run_service.airspace_daily_history.location
  project     = google_cloud_run_service.airspace_daily_history.project
  service     = google_cloud_run_service.airspace_daily_history.name
  policy_data = data.google_iam_policy.noauth.policy_data
}

# # # Aircraft History service
resource "google_cloud_run_service" "airspace_monthly_history" {
  name     = "airspace-monthly-history"
  location = var.region

  template {
    spec {
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/airspace_monthly_history:latest"
        env {
          name  = "GOOGLE_APPLICATION_CREDENTIALS"
          value = var.google_list_credentials
        }
        env {
          name  = "GOOGLE_CLOUD_PROJECT_ID"
          value = var.project_id
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
  location    = google_cloud_run_service.airspace_monthly_history.location
  project     = google_cloud_run_service.airspace_monthly_history.project
  service     = google_cloud_run_service.airspace_monthly_history.name
  policy_data = data.google_iam_policy.noauth.policy_data
}

# # # aircraft positions

resource "google_cloud_run_service" "aircraft_positions" {
  name     = "aircraft-positions"
  location = var.region

  template {
    spec {
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/aircraft_positions:latest"
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

resource "google_cloud_run_service_iam_policy" "noauth_aircraft_positions" {
  location = google_cloud_run_service.aircraft_positions.location
  project  = google_cloud_run_service.aircraft_positions.project
  service  = google_cloud_run_service.aircraft_positions.name

  policy_data = data.google_iam_policy.noauth.policy_data
}