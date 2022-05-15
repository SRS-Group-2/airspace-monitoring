# # using deployment instead of Job because we don't expect it to end (it's a continuous job for us)
# resource "kubernetes_deployment" "flink_jobmanager" {
#   depends_on = [kubernetes_namespace.main_namespace]
#   metadata {
#     name      = "flink-jobmanager"
#     namespace = var.kube_namespace
#   }

#   spec {
#     replicas = 1

#     selector {
#       match_labels = {
#         app = "flink"

#         component = "jobmanager"
#       }
#     }

#     template {
#       metadata {
#         labels = {
#           app = "flink"

#           component = "jobmanager"
#         }
#       }

#       spec {
#         service_account_name = kubernetes_service_account.flink_kube_account.metadata[0].name
#         volume {
#           name = "flink-config-volume"

#           config_map {
#             name = "flink-config"

#             items {
#               key  = "flink-conf.yaml"
#               path = "flink-conf.yaml"
#             }

#             items {
#               key  = "log4j-console.properties"
#               path = "log4j-console.properties"
#             }
#           }
#         }

#          init_container {
#            name  = "workload-identity-initcontainer"
#            image = "alpine/curl:3.14" // gcr.io/google.com/cloudsdktool/cloud-sdk:326.0.0-alpine
#            command = [
#              "/bin/sh",
#              "-c",
#              "echo Going to sleep it out && sleep 20 && (curl -s -H 'Metadata-Flavor: Google' 'http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token' --retry 30 --retry-connrefused --retry-max-time 30 > /dev/null && echo Metadata server working) || exit 1"
#            ]
#          }

#         container {
#           name  = "jobmanager"
#           image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/states_source:latest"
#           args  = ["standalone-job", "--job-classname", "it.unibo.states_source.Main"]

#           env {
#             name  = "GOOGLE_CLOUD_PROJECT_ID"
#             value = var.project_id
#           }

#           env {
#             name  = "GOOGLE_PUBSUB_VECTORS_TOPIC_ID"
#             value = var.vectors_topic
#           }

#           env {
#             name  = "COORDINATES"
#             value = var.opensky_bb
#           }

#           port {
#             name           = "rpc"
#             container_port = 6123
#           }

#           port {
#             name           = "blob-server"
#             container_port = 6124
#           }

#           port {
#             name           = "webui"
#             container_port = 8081
#           }

#           volume_mount {
#             name       = "flink-config-volume"
#             mount_path = "/opt/flink/conf"
#           }

#           liveness_probe {
#             tcp_socket {
#               port = "6123"
#             }

#             initial_delay_seconds = 60
#             period_seconds        = 60
#           }

#           security_context {
#             run_as_user = 9999
#           }
#         }

#         node_selector = {
#           "iam.gke.io/gke-metadata-server-enabled" = "true"
#         }
#       }
#     }
#   }
# }
