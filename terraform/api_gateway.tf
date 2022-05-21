
locals {
  api_config_id_prefix     = "airspace-monitoring"
  api_gateway_container_id = "api-gw"
  gateway_id               = "airspace-monitoring"
}


data "template_file" "aircraft_info" {
  template = file("airspace_monitoring.yaml")
  vars = {
    aircraft_info_url            = "${google_cloud_run_service.aircraft_info.status[0].url}"
    aircraft_list_url            = "${google_cloud_run_service.aircraft_list.status[0].url}"
    airspace_daily_history_url   = "${google_cloud_run_service.airspace_daily_history.status[0].url}"
    airspace_monthly_history_url = "${google_cloud_run_service.airspace_monthly_history.status[0].url}"
    airspace_web_ui              = "${google_cloud_run_service.web_ui.status[0].url}"
  }
  depends_on = [google_cloud_run_service.aircraft_info, google_cloud_run_service.aircraft_list, google_cloud_run_service.airspace_daily_history, google_cloud_run_service.airspace_monthly_history, google_cloud_run_service.web_ui]
}


resource "google_api_gateway_api" "api_gw" {
  provider     = google-beta
  api_id       = local.api_gateway_container_id
  display_name = "airspace-info"
}


resource "google_api_gateway_api_config" "api_cfg" {
  provider             = google-beta
  api                  = google_api_gateway_api.api_gw.api_id
  api_config_id_prefix = local.api_config_id_prefix
  display_name         = "airspace-info"

  openapi_documents {
    document {
      path     = "aircraft_info.yaml"
      contents = base64encode(data.template_file.aircraft_info.rendered)

    }

  }
}

resource "google_api_gateway_gateway" "gw" {
  provider = google-beta
  region   = var.region

  api_config = google_api_gateway_api_config.api_cfg.id

  gateway_id   = local.gateway_id
  display_name = "The Gateway for airspace-monitoring"

  depends_on = [google_api_gateway_api_config.api_cfg]
}

