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
  default_max_pods_per_node  = 32 # using a /24 subnet per pool, so actually a /26 per node, considering 4 nodes, so https://cloud.google.com/kubernetes-engine/docs/how-to/flexible-pod-cidr#cidr_ranges_for_clusters says max 32

  node_pools = [
    {
      name               = "micro-node-pool"
      machine_type       = "e2-micro"
      node_locations     = "${var.region}-c"
      min_count          = 1
      max_count          = 3
      local_ssd_count    = 0
      disk_size_gb       = 10
      disk_type          = "pd-standard"
      image_type         = "COS_CONTAINERD"
      auto_repair        = true
      auto_upgrade       = true
      preemptible        = false
      initial_node_count = 1
    },
    {
      name               = "small-node-pool"
      machine_type       = "e2-small"
      node_locations     = "${var.region}-c"
      min_count          = 1
      max_count          = 2
      local_ssd_count    = 0
      disk_size_gb       = 10
      disk_type          = "pd-standard"
      image_type         = "COS_CONTAINERD"
      auto_repair        = true
      auto_upgrade       = true
      preemptible        = false
      initial_node_count = 1
    },
    {
      name               = "medium-node-pool"
      machine_type       = "e2-medium"
      node_locations     = "${var.region}-c"
      min_count          = 1
      max_count          = 1
      local_ssd_count    = 0
      disk_size_gb       = 10
      disk_type          = "pd-standard"
      image_type         = "COS_CONTAINERD"
      auto_repair        = true
      auto_upgrade       = true
      preemptible        = false
      initial_node_count = 1
    },
  ]

  node_pools_oauth_scopes = {
    all = [
      "https://www.googleapis.com/auth/devstorage.read_only", # to read from artifact registry
    ]

    micro-node-pool = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/devstorage.read_only", # to read from artifact registry
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/service.management.readonly",
      "https://www.googleapis.com/auth/servicecontrol",
      "https://www.googleapis.com/auth/trace.append",
      "https://www.googleapis.com/auth/datastore", # is this necessary?
    ]

    small-node-pool = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/devstorage.read_only", # to read from artifact registry
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/service.management.readonly",
      "https://www.googleapis.com/auth/servicecontrol",
      "https://www.googleapis.com/auth/trace.append",
      "https://www.googleapis.com/auth/datastore", # is this necessary?
    ]

    medium-node-pool = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/devstorage.read_only", # to read from artifact registry
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/service.management.readonly",
      "https://www.googleapis.com/auth/servicecontrol",
      "https://www.googleapis.com/auth/trace.append",
      "https://www.googleapis.com/auth/datastore", # is this necessary?
    ]
  }

  node_pools_labels = {
    all = {}

    micro-node-pool = {
      node_type = "micro"
    }

    small-node-pool = {
      node_type  = "small"
    }

    medium-node-pool = {
      node_type  = "medium"
    }
  }
}

resource "google_project_iam_member" "allow_image_pull" {
  project = var.project_id
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:${module.gke.service_account}"
}
