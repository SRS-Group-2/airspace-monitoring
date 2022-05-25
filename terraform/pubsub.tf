# Google PubSub
resource "google_pubsub_topic" "pubsub_positions" {
  depends_on = [
    google_project_service.pubsub,
  ]
  project = var.project_id
  name    = var.vectors_topic
}
