# general
variable "project_id" {
  type        = string
  description = "The Google Cloud Project Id"
}

variable "region" {
  type        = string
  description = "The region where to deploy"
}

variable "google_list_credentials" {
  type        = string
  description = "gcp credentials"
}


variable "vectors_topic" {
  type        = string
  description = "The name of the PubSub topic with the state vectors"
}

variable "history_1_topic" {
  type        = string
  description = "The name of the PubSub topic with the distance and CO2 of the last hour"
}

variable "history_6_topic" {
  type        = string
  description = "The name of the PubSub topic with the distance and CO2 of the last six hours"
}

variable "history_24_topic" {
  type        = string
  description = "The name of the PubSub topic with the distanche and CO2 of the last twentyfour hours"
}


variable "google_monthly_history_credentials" {
  type        = string
  description = "The credentials for the airspace monthly history service"
}

# kubernetes
# variable "kube_cluster" {
#   type        = string
#   description = "Name of the kubernetes cluster being deployed"
# }

# variable "kube_network" {
#   type        = string
#   description = "Name of the VPC network"
# }

# variable "kube_subnetwork" {
#   type        = string
#   description = "Name of the VPC subnetwork"
# }

# variable "kube_pods_range" {
#   type        = string
#   description = "Name of ip range of pods in kube subnetwork"
# }

# variable "kube_services_range" {
#   type        = string
#   description = "Name of ip range of services in kube subnetwork"
# }

# variable "flink_version" {
#   type        = string
#   description = "Label of Flink container image"
# }
