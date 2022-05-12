# resource "kubernetes_deployment" "airspace_history_calculator" {
#   metadata {
#     name = "airspace-history-calculator"
#   }

#   spec {
#     replicas = 1

#     selector {
#       match_labels = {
#         app = "history-calculator"
#       }
#     }

#     template {
#       metadata {
#         labels = {
#           app = "history-calculator"
#         }
#       }

#       spec {
#         container {
#           name  = "airspace_history_calculator"
#           image = "${var.region}-docker.pkg.dev/${var.project_id}/docker-repo/airspace_history_calculator:latest"

#           env {
#             name  = "GOOGLE_APPLICATION_CREDENTIALS"
#             value = var.airspace_history_calculator_credentials
#           }
#           env {
#             name  = "GOOGLE_CLOUD_PROJECT_ID"
#             value = var.project_id
#           }
#           env {
#             name  = "GIN_MODE"
#             value = "release"
#           }

#           #           port {
#           #             name           = "port"
#           #             container_port = 8080
#           #           }

#           #           liveness_probe {
#           #             tcp_socket {
#           #               port = "6122"
#           #             }

#           #             initial_delay_seconds = 60
#           #             period_seconds        = 60
#           #           }

#           security_context {
#             run_as_user = 9999
#           }
#         }

#       }
#     }
#   }
# }
