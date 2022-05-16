# general
variable "project_id" {
  type        = string
  description = "The Google Cloud Project Id"
}

variable "region" {
  type        = string
  description = "The region where to deploy"
}

# pubsub
variable "vectors_topic" {
  type        = string
  description = "The name of the PubSub topic with the state vectors"
}

# kubernetes
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
