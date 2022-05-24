# general
variable "project_id" {
  type        = string
  description = "The Google Cloud Project Id"
}

variable "region" {
  type        = string
  description = "The region where to deploy"
}

# docker
variable "docker_repo_name" {
  type        = string
  description = "The name of the Google Docker repository in which the images are stored"
}

variable "docker_repo_region" {
  type        = string
  description = "The region of the Google Docker repository in which the images are stored"
}

variable "aircraft_list_tag" {
  type        = string
  description = "The tag for the aircraft_list Docker image"
}
variable "airspace_daily_history_tag" {
  type        = string
  description = "The tag for the airspace_daily_history Docker image"
}
variable "websocket_endpoints_tag" {
  type        = string
  description = "The tag for the websocket_endpoints Docker image"
}
variable "airspace_history_calculator_tag" {
  type        = string
  description = "The tag for the airspace_history_calculator Docker image"
}
variable "airspace_monthly_history_tag" {
  type        = string
  description = "The tag for the airspace_monthly_history Docker image"
}
variable "aircraft_positions_tag" {
  type        = string
  description = "The tag for the aircraft_positions Docker image"
}
variable "aircraft_info_tag" {
  type        = string
  description = "The tag for the aircraft_info Docker image"
}
variable "web_ui_tag" {
  type        = string
  description = "The tag for the web_ui Docker image"
}
variable "states_source_tag" {
  type        = string
  description = "The tag for the states_source Docker image"
}

# pubsub
variable "vectors_topic" {
  type        = string
  description = "The name of the PubSub topic with the state vectors"
}

#kubernetes
variable "kube_cluster" {
  type        = string
  description = "Name of the kubernetes cluster being deployed"
}

variable "kube_network" {
  type        = string
  description = "Name of the VPC network"
}

variable "kube_subnetwork" {
  type        = string
  description = "Name of the VPC subnetwork"
}

variable "kube_pods_range" {
  type        = string
  description = "Name of ip range of pods in kube subnetwork"
}

variable "kube_services_range" {
  type        = string
  description = "Name of ip range of services in kube subnetwork"
}

variable "kube_namespace" {
  type        = string
  description = "Name of the namespace in which the applications will be deployed"
}

# app
variable "opensky_bb" {
  type        = string
  description = "Env for flink application"
}
