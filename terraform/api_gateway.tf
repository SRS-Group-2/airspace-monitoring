
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
  }
  depends_on = [google_cloud_run_service.aircraft_info, google_cloud_run_service.aircraft_list, google_cloud_run_service.airspace_daily_history, google_cloud_run_service.airspace_monthly_history]
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



///////////////////////////////////////////////////////////////////////////////////////////////////////////



# resource "google_api_gateway_api" "api_gw2" {
#   provider     = google-beta
#   api_id       = local.api_gateway_container_id2
#   display_name = "airspace-list"

# }



# resource "google_api_gateway_api_config" "api_cfg2" {
#   provider             = google-beta
#   api                  = google_api_gateway_api.api_gw2.api_id
#   api_config_id_prefix = local.api_config_id_prefix
#   display_name         = "airspace-list"

#   openapi_documents {
#     document {
#       path     = "aircraft_list.yaml"
#       contents = filebase64("aircraft_list.yaml")
#     }
#   }

# }


# resource "google_api_gateway_gateway" "gw2" {
#   provider = google-beta
#   region   = var.region

#   api_config = google_api_gateway_api_config.api_cfg2.id

#   gateway_id   = local.gateway_id2
#   display_name = "The Gateway for airspace-monitoring2"

#   depends_on = [google_api_gateway_api_config.api_cfg2]
# }



///////////////////////////////////////////


# resource "google_api_gateway_api" "api_gw3" {
#   provider     = google-beta
#   api_id       = local.api_gateway_container_id3
#   display_name = "airspace-daily-history"

# }



# resource "google_api_gateway_api_config" "api_cfg3" {
#   provider             = google-beta
#   api                  = google_api_gateway_api.api_gw3.api_id
#   api_config_id_prefix = local.api_config_id_prefix
#   display_name         = "airspace-daily-history"

#   openapi_documents {
#     document {
#       path     = "airspace-daily-history.yaml"
#       contents = filebase64("airspace-daily-history.yaml")
#     }
#   }

# }


# resource "google_api_gateway_gateway" "gw3" {
#   provider = google-beta
#   region   = var.region

#   api_config = google_api_gateway_api_config.api_cfg3.id

#   gateway_id   = local.gateway_id3
#   display_name = "The Gateway for daili history"

#   depends_on = [google_api_gateway_api_config.api_cfg3]
# }


///////////////////////////////////////////////////////////////

# resource "google_api_gateway_api" "api_gw4" {
#   provider     = google-beta
#   api_id       = local.api_gateway_container_id4
#   display_name = "airspace-daily-history"

# }



# resource "google_api_gateway_api_config" "api_cfg4" {
#   provider             = google-beta
#   api                  = google_api_gateway_api.api_gw4.api_id
#   api_config_id_prefix = local.api_config_id_prefix
#   display_name         = "airspace-monthly-history"

#   openapi_documents {
#     document {
#       path     = "airspace-monthly-history.yaml"
#       contents = filebase64("airspace-monthly-history.yaml")
#     }
#   }

# }


# resource "google_api_gateway_gateway" "gw4" {
#   provider = google-beta
#   region   = var.region

#   api_config = google_api_gateway_api_config.api_cfg4.id

#   gateway_id   = local.gateway_id4
#   display_name = "The Gateway for monthly history"

#   depends_on = [google_api_gateway_api_config.api_cfg4]
# }

# /////

