resource "kubernetes_deployment" "flink_taskmanager" {
  metadata {
    name = "flink-taskmanager"
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

        container {
          name  = "taskmanager"
          image = "us-central1-docker.pkg.dev/master-choir-347215/docker-repo/flink_quickstart:latest"
          args  = ["taskmanager"]

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

            initial_delay_seconds = 30
            period_seconds        = 60
          }

          security_context {
            run_as_user = 9999
          }
        }
      }
    }
  }
}
