# resource "kubernetes_service" "flink_jobmanager" {
#   depends_on = [kubernetes_namespace.main_namespace]
#   metadata {
#     name      = "flink-jobmanager"
#     namespace = var.kube_namespace
#   }

#   spec {
#     port {
#       name = "rpc"
#       port = 6123
#     }

#     port {
#       name = "blob-server"
#       port = 6124
#     }

#     port {
#       name = "webui"
#       port = 8081
#     }

#     selector = {
#       app = "flink"

#       component = "jobmanager"
#     }

#     type = "ClusterIP"
#   }
# }
