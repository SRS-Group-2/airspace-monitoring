# Google PubSub
/*resource "google_pubsub_topic" "pubsub_list" {
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

# resource "google_pubsub_subscription" "pubsub_sub_list" {
#   name  = var.aircraft_list_subscriber_id
#   topic = google_pubsub_topic.pubsub_list.name

#   # 20 minutes
#   message_retention_duration = "1200s"
#   retain_acked_messages      = true

#   ack_deadline_seconds = 20

#   expiration_policy {
#     ttl = "300000.5s"
#   }
# }
*/