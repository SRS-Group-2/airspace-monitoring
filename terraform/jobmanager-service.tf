resource "kubernetes_service" "flink_jobmanager" {
  metadata {
    name = "flink-jobmanager"
  }

  spec {
    port {
      name = "rpc"
      port = 6123
    }

    port {
      name = "blob-server"
      port = 6124
    }

    port {
      name = "webui"
      port = 8081
    }

    selector = {
      app = "flink"

      component = "jobmanager"
    }

    type = "ClusterIP"
  }
}
