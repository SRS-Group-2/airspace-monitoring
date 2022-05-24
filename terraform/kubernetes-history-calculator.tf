# # service account for airspace history calculator
# resource "google_service_account" "airspace_history_calculator_sa" {
#   account_id   = "airspace-history-calculator"
#   display_name = "A service account for the Airspace History Calculator"
# }

# # bind the service account to the necessary roles
# resource "google_project_iam_binding" "airspace_history_calculator_binding" {
#   project = var.project_id
#   role    = "roles/datastore.user"

#   members = [
#     "serviceAccount:${google_service_account.airspace_history_calculator_sa.email}",
#   ]
# }

locals {
  # airspace_history_calculator_name = google_service_account.airspace_history_calculator.name
  # airspace_history_calculator_email = google_service_account.airspace_history_calculator.email
  airspace_history_calculator_name  = "projects/${var.project_id}/serviceAccounts/airspace-history-calculator@${var.project_id}.iam.gserviceaccount.com"
  airspace_history_calculator_email = "airspace-history-calculator@${var.project_id}.iam.gserviceaccount.com"
}

# resource "google_service_account_key" "airspace_history_calculator_key" {
#   service_account_id = local.airspace_history_calculator_name
#   public_key_type    = "TYPE_X509_PEM_FILE"
# }

# kubernetes service account for airspace history calculator
resource "kubernetes_service_account" "airspace_history_calculator_kube_account" {
  depends_on = [
    kubernetes_namespace.main_namespace,
    kubernetes_deployment.flink_jobmanager,  # requires the data from these two to work
    kubernetes_deployment.flink_taskmanager, # requires the data from these two to work
    #google_project_iam_binding.airspace_daily_history_binding_log
  ]
  metadata {
    name      = "airspace-history-calculator-account"
    namespace = var.kube_namespace
    annotations = {
      "iam.gke.io/gcp-service-account" = local.airspace_history_calculator_email
    }
  }
}

# bind service account and kubernetes service account
resource "google_service_account_iam_binding" "airspace_history_calculator_accounts_binding" {
  service_account_id = local.airspace_history_calculator_name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:${var.project_id}.svc.id.goog[${var.kube_namespace}/${kubernetes_service_account.airspace_history_calculator_kube_account.metadata[0].name}]",
  ]
}

resource "kubernetes_deployment" "airspace_history_calculator" {
  depends_on = [
    kubernetes_namespace.main_namespace,
    # google_service_account.airspace_history_calculator_sa,
    # google_project_iam_binding.airspace_history_calculator_binding,
    # google_service_account_key.airspace_history_calculator_key,
  ]
  metadata {
    name      = "airspace-history-calculator"
    namespace = var.kube_namespace
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        "app" = "history-calculator"
      }
    }

    template {
      metadata {
        labels = {
          app = "history-calculator"
        }
      }

      spec {
        service_account_name = kubernetes_service_account.airspace_history_calculator_kube_account.metadata[0].name
        init_container {
          name  = "workload-identity-initcontainer"
          image = "alpine/curl:3.14" // "gcr.io/google.com/cloudsdktool/cloud-sdk:385.0.0-alpine" //  
          command = [
            "/bin/sh",
            "-c",
            "echo Going to sleep it out && sleep 20 && (curl -s -H 'Metadata-Flavor: Google' 'http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token' --retry 30 --retry-connrefused --retry-max-time 30 > /dev/null && echo Metadata server working) || exit 1"
          ]
        }
        container {
          name  = "airspace-history-calculator"
          image = "${var.docker_repo_region}-docker.pkg.dev/${var.project_id}/${var.docker_repo_name}/airspace_history_calculator:${var.airspace_history_calculator_tag}"

          env {
            name  = "GOOGLE_CLOUD_PROJECT_ID"
            value = var.project_id
          }
          env {
            name  = "GIN_MODE"
            value = "release"
          }

          security_context {
            run_as_user = 9999
          }
        }

        node_selector = {
          "iam.gke.io/gke-metadata-server-enabled" = "true"
        }
      }
    }
  }
}
