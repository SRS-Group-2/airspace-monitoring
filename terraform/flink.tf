
# Flink basic cluster
resource "google_cloud_run_service" "flink_job_manager" {
  name     = "job-manager"
  location = var.region

  template {
    spec {
      containers {
        image   = "flink:${var.flink_version}"
        command = ["sh", "-c", "jobmanager"]
        env {
          name  = "FLINK_PROPERTIES"
          value = "jobmanager.rpc.address: https://job-manager-${var.project_id_hash}.a.run.app"
        }
      }
    }
  }
}

resource "google_cloud_run_service" "flink_task_manager" {
  name     = "task-manager"
  location = var.region

  template {
    spec {
      containers {
        image   = "flink:${var.flink_version}"
        command = ["sh", "-c", "taskmanager"]
        env {
          name  = "FLINK_PROPERTIES"
          value = "jobmanager.rpc.address: https://job-manager-${var.project_id_hash}.a.run.app"
        }
      }
    }
  }
}
