locals {
  api_config_id_prefix      = "api"
  ////
  api_gateway_container_id  = "api-gw"
  api_gateway_container_id2 = "api-gw2"
  api_gateway_container_id3 = "api-gw3"
  ////
  gateway_id                = "gw"
  gateway_id2               = "gw2"
  gateway_id3               = "gw3"

}

///////////////////GOOGLE_API_GATEWAY////////////////////////
//1
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
      contents = filebase64("aircraft_info.yaml")
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




resource "google_api_gateway_api" "api_gw2" {
  provider     = google-beta
  api_id       = local.api_gateway_container_id2
  display_name = "airspace-list"

}



resource "google_api_gateway_api_config" "api_cfg2" {
  provider             = google-beta
  api                  = google_api_gateway_api.api_gw2.api_id
  api_config_id_prefix = local.api_config_id_prefix
  display_name         = "airspace-list"

  openapi_documents {
    document {
      path     = "aircraft_list.yaml"
      contents = filebase64("aircraft_list.yaml")
    }
  }

}


resource "google_api_gateway_gateway" "gw2" {
  provider = google-beta
  region   = var.region

  api_config = google_api_gateway_api_config.api_cfg2.id

  gateway_id   = local.gateway_id2
  display_name = "The Gateway for airspace-monitoring2"

  depends_on = [google_api_gateway_api_config.api_cfg2]
}


    
///////////////////////////////////////////


resource "google_api_gateway_api" "api_gw3" {
  provider     = google-beta
  api_id       = local.api_gateway_container_id3
  display_name = "airspace-daily-history"

}



resource "google_api_gateway_api_config" "api_cfg3" {
  provider             = google-beta
  api                  = google_api_gateway_api.api_gw3.api_id
  api_config_id_prefix = local.api_config_id_prefix
  display_name         = "airspace-daily-history"

  openapi_documents {
    document {
      path     = "airspace-daily-historyt.yaml"
      contents = filebase64("airspace-daily-history.yaml")
    }
  }

}


resource "google_api_gateway_gateway" "gw3" {
  provider = google-beta
  region   = var.region

  api_config = google_api_gateway_api_config.api_cfg3.id

  gateway_id   = local.gateway_id3
  display_name = "The Gateway for daili history"

  depends_on = [google_api_gateway_api_config.api_cfg3]
}