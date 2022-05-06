resource "kubernetes_job" "flink_jobmanager" {
  metadata {
    name = "flink-jobmanager"
  }

  spec {
    template {
      metadata {
        labels = {
          app = "flink"

          component = "jobmanager"
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
          name  = "jobmanager"
          image = "us-central1-docker.pkg.dev/master-choir-347215/docker-repo/flink_quickstart:latest"
          args  = ["standalone-job", "--job-classname", "org.myorg.quickstart.WordCount"]

          port {
            name           = "rpc"
            container_port = 6123
          }

          port {
            name           = "blob-server"
            container_port = 6124
          }

          port {
            name           = "webui"
            container_port = 8081
          }

          volume_mount {
            name       = "flink-config-volume"
            mount_path = "/opt/flink/conf"
          }

          liveness_probe {
            tcp_socket {
              port = "6123"
            }

            initial_delay_seconds = 30
            period_seconds        = 60
          }

          security_context {
            run_as_user = 9999
          }
        }

        restart_policy = "OnFailure"
      }
    }
  }
}
