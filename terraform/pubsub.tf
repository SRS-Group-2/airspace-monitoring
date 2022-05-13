# Google PubSub
resource "google_pubsub_topic" "pubsub_positions" {
  project = var.project_id
  name    = var.vectors_topic
}
