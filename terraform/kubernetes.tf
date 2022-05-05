data "google_client_config" "default" {}

provider "kubernetes" {
  host                   = "https://${module.gke.endpoint}"
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(module.gke.ca_certificate)
}

module "gke" {
  depends_on = [google_compute_subnetwork.kubernetes_subnetwork]

  source                     = "terraform-google-modules/kubernetes-engine/google"
  project_id                 = var.project_id
  name                       = var.kube_cluster
  region                     = var.region
  zones                      = ["${var.region}-c"]
  network                    = var.kube_network
  subnetwork                 = var.kube_subnetwork
  ip_range_pods              = var.kube_pods_range
  ip_range_services          = var.kube_services_range
  http_load_balancing        = false
  network_policy             = false
  horizontal_pod_autoscaling = true
  filestore_csi_driver       = false
  default_max_pods_per_node  = 32 # using a /24 subnet per pool, so actually a /26 per node, so https://cloud.google.com/kubernetes-engine/docs/how-to/flexible-pod-cidr#cidr_ranges_for_clusters says max 32

  node_pools = [
    {
      name               = "default-node-pool"
      machine_type       = "e2-micro"
      node_locations     = "${var.region}-c"
      min_count          = 1
      max_count          = 4
      local_ssd_count    = 0
      disk_size_gb       = 32
      disk_type          = "pd-standard"
      image_type         = "COS_CONTAINERD"
      auto_repair        = true
      auto_upgrade       = true
      preemptible        = false
      initial_node_count = 3
    },
  ]
}
