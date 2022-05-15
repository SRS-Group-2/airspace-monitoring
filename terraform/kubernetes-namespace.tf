resource "kubernetes_namespace" "main_namespace" {
  depends_on = [module.gke]
  metadata {
    name = var.kube_namespace
  }
}
