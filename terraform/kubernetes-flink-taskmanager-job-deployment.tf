# resource "kubernetes_deployment" "flink_taskmanager" {
#   metadata {
#     name      = "flink-taskmanager"
#     namespace = var.kube_namespace
#   }

#   spec {
#     replicas = 1

#     selector {
#       match_labels = {
#         app = "flink"

#         component = "taskmanager"
#       }
#     }

#     template {
#       metadata {
#         labels = {
#           app = "flink"

#           component = "taskmanager"
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
#           name  = "taskmanager"
#           image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/states_source:latest"
#           args  = ["taskmanager"]

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
#             name  = "COORDINATES"
#             value = var.opensky_bb
#           }

#           port {
#             name           = "rpc"
#             container_port = 6122
#           }

#           port {
#             name           = "query-state"
#             container_port = 6125
#           }

#           volume_mount {
#             name       = "flink-config-volume"
#             mount_path = "/opt/flink/conf/"
#           }

#           liveness_probe {
#             tcp_socket {
#               port = "6122"
#             }

#             initial_delay_seconds = 60
#             period_seconds        = 60
#           }

#           security_context {
#             run_as_user = 9999
#           }
#         }

#         node_selector = {
#           node_type = "small"
#         }
#       }
#     }
#   }
# }
