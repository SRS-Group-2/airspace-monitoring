variable "project_id" {
  type        = string
  description = "The Google Cloud Project Id"
}

variable "project_id_hash" {
  type        = string
  description = "The hash of the project id"
}

variable "region" {
  type        = string
  description = "The region where to deploy"
}

variable "aircraft_list_topic" {
  type        = string
  description = "The name of the PubSub topic with the list of aircrafts"
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

variable "aircraft_list_subscriber_id" {
  type        = string
  description = "The id of the subscriber for the aircraft_list"
}

variable "google_list_credentials" {
  type        = string
  description = "The credentials for the aircraft list service"
}

variable "flink_version" {
  type        = string
  description = "Version of the Flink instances being deployed"
}
