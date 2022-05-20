# define Terraform general information
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.20.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.11.0"
    }
  }

  # where to save states
  backend "gcs" {
    bucket = var.state_bucket
    prefix = "terraform/state"
  }
}

# define cloud provider
provider "google" {
  project = var.project_id
  region  = var.region
  zone    = "${var.region}-c"
}
