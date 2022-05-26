resource "google_compute_subnetwork" "kubernetes_subnetwork" {
  name          = var.kube_subnetwork
  ip_cidr_range = "10.2.0.0/16"
  region        = var.region
  network       = google_compute_network.kubernetes_network.id
  project       = var.project_id
  secondary_ip_range {
    range_name    = var.kube_pods_range
    ip_cidr_range = "192.168.10.0/24"
  }
  secondary_ip_range {
    range_name    = var.kube_services_range
    ip_cidr_range = "192.168.11.0/24"
  }
}

resource "google_compute_network" "kubernetes_network" {
  name                    = var.kube_network
  auto_create_subnetworks = false
}
