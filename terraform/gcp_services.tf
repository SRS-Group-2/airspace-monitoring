# Google PubSub
resource "google_pubsub_topic" "pubsub_list" {
  project = var.project_id
  name    = var.aircraft_list_topic
}

resource "google_pubsub_topic" "pubsub_positions" {
  project = var.project_id
  name    = var.vectors_topic
}

resource "google_pubsub_topic" "pubsub_1h" {
  project = var.project_id
  name    = var.history_1_topic
}

resource "google_pubsub_topic" "pubsub_6h" {
  project = var.project_id
  name    = var.history_6_topic
}

resource "google_pubsub_topic" "pubsub_24h" {
  project = var.project_id
  name    = var.history_24_topic
}

# Subscriptions will be created by containers when necessary. 
# As per Google policy, they are deleted after 31 days of inactivity
