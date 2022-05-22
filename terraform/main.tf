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
    bucket = "choir-tf-state"
    prefix = "terraform/state"
  }
}

# define cloud provider
provider "google" {
  project = var.project_id
  region  = var.region
  zone    = "${var.region}-c"
}

provider "google-beta" {
  project = var.project_id
  region  = var.region
  zone    = "${var.region}-c"
}

data "google_project" "project" {
}

