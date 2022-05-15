# service account for flink
resource "google_service_account" "flink" {
  account_id   = "flink-source"
  display_name = "A service account for Flink"
}

# bind the service account to the necessary roles
resource "google_project_iam_binding" "flink_firestore_binding_user" {
  project = var.project_id
  role    = "roles/datastore.user"

  members = [
    "serviceAccount:${google_service_account.flink.email}",
  ]
}

resource "google_project_iam_binding" "flink_pubsub_binding" {
  project = var.project_id
  role    = "roles/pubsub.publisher" // TODO check whether editor is necessary

  members = [
    "serviceAccount:${google_service_account.flink.email}",
  ]
}

# kubernetes service account for airspace history calculator
resource "kubernetes_service_account" "flink_kube_account" {
  depends_on = [kubernetes_namespace.main_namespace]
  metadata {
    name      = "flink-account"
    namespace = var.kube_namespace
    annotations = {
      "iam.gke.io/gcp-service-account" = google_service_account.flink.email
    }
  }
}

# bind service account and kubernetes service account
resource "google_service_account_iam_binding" "flink_accounts_binding" {
  service_account_id = google_service_account.flink.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:${var.project_id}.svc.id.goog[${var.kube_namespace}/${kubernetes_service_account.flink_kube_account.metadata[0].name}]",
  ]
}

resource "kubernetes_deployment" "flink_taskmanager" {
  depends_on = [kubernetes_namespace.main_namespace]
  metadata {
    name      = "flink-taskmanager"
    namespace = var.kube_namespace
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "flink"

        component = "taskmanager"
      }
    }

    template {
      metadata {
        labels = {
          app = "flink"

          component = "taskmanager"
        }
      }

      spec {
        service_account_name = kubernetes_service_account.flink_kube_account.metadata[0].name
        volume {
          name = "flink-config-volume"

          config_map {
            name = "flink-config"

            items {
              key  = "flink-conf.yaml"
              path = "flink-conf.yaml"
            }

            items {
              key  = "log4j-console.properties"
              path = "log4j-console.properties"
            }
          }
        }

         init_container {
           name  = "workload-identity-initcontainer"
           image = "alpine/curl:3.14" // gcr.io/google.com/cloudsdktool/cloud-sdk:326.0.0-alpine
           command = [
             "/bin/sh",
             "-c",
             "echo Going to sleep it out && sleep 30 && (curl -s -H 'Metadata-Flavor: Google' 'http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token' --retry 30 --retry-connrefused --retry-max-time 30 > /dev/null && echo Metadata server working) || exit 1"
           ]
         }

        container {
          name  = "taskmanager"
          image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/states_source:latest"
          args  = ["taskmanager"]

          env {
            name  = "GOOGLE_CLOUD_PROJECT_ID"
            value = var.project_id
          }

          env {
            name  = "GOOGLE_PUBSUB_VECTORS_TOPIC_ID"
            value = var.vectors_topic
          }

          env {
            name  = "COORDINATES"
            value = var.opensky_bb
          }

          port {
            name           = "rpc"
            container_port = 6122
          }

          port {
            name           = "query-state"
            container_port = 6125
          }

          volume_mount {
            name       = "flink-config-volume"
            mount_path = "/opt/flink/conf/"
          }

          liveness_probe {
            tcp_socket {
              port = "6122"
            }

            initial_delay_seconds = 60
            period_seconds        = 60
          }

          security_context {
            run_as_user = 9999
          }
        }

        node_selector = {
          node_type                                = "small"
          "iam.gke.io/gke-metadata-server-enabled" = "true"
        }
      }
    }
  }
}
