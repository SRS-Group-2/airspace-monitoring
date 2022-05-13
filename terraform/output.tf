output "service_url" {
  value = google_cloud_run_service.aircraft_info.status[0].url
}