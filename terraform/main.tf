# define Terraform general information
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.5.0"
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
