# # using deployment instead of Job because we don't expect it to end (it's a continuous job for us)
# resource "kubernetes_deployment" "flink_jobmanager" {
#   metadata {
#     name = "flink-jobmanager"
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

#         container {
#           name  = "jobmanager"
#           image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/states_source:latest"
#           args  = ["standalone-job", "--job-classname", "it.unibo.states_source.Main"]

#           env {
#             name  = "GOOGLE_APPLICATION_CREDENTIALS"
#             value = var.states_source_credentials
#           }

#           env {
#             name  = "GOOGLE_CLOUD_PROJECT_ID"
#             value = var.project_id
#           }

#           env {
#             name  = "GOOGLE_PUBSUB_VECTORS_TOPIC_ID"
#             value = var.vectors_topic
#           }

#           env {
#             name   = "COORDINATES"
#             values = var.opensky_bb
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
#       }
#     }
#   }
# }
