const base_aircrafts_route = "/aircrafts"
const base_distance_route = "/distance"
const base_co2_route = "/co2"

const max_days_range = 30

function request_value(url, value_callback) {
  var xhttp = new XMLHttpRequest()
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      value_callback(xhttp.responseText)
    }
  }
  xhttp.open("GET", url)
  xhttp.send()
}

function request_value_with_params(url, value_callback, params) {
  var newUrl = url + "?"
  var keys = Object.keys(params)
  for (const key in keys) {
    newUrl = newUrl + keys[key] + "=" + params[keys[key]] + "&"
  }
  newUrl = newUrl.slice(0, -1)
  console.log(newUrl)
  request_value(newUrl, value_callback)
}

function add_time_selectors() {
  const begin_time_selector_label = document.createElement("label")
  begin_time_selector_label.for = "begin_time_selector"
  document.getElementById("historic_begin_section").append(begin_time_selector_label)
  const begin_time_selector = document.createElement("input")
  begin_time_selector.type = "time"
  begin_time_selector.required = true
  begin_time_selector.id = "begin_time_selector"
  begin_time_selector.name = "begin_time_selector"
  document.getElementById("historic_begin_section").append(begin_time_selector)
  const end_time_selector_label = document.createElement("label")
  end_time_selector_label.for = "end_time_selector"
  document.getElementById("historic_end_section").append(end_time_selector_label)
  const end_time_selector = document.createElement("input")
  end_time_selector.type = "time"
  end_time_selector.required = true
  end_time_selector.id = "end_time_selector"
  end_time_selector.name = "end_time_selector"
  document.getElementById("historic_end_section").append(end_time_selector)
}

function get_begin_timestamp() {
  return get_timestamp(document.getElementById("historic_begin_date"), document.getElementById("begin_time_selector"))
}

function get_end_timestamp() {
  return get_timestamp(document.getElementById("historic_end_date"), document.getElementById("end_time_selector"))
}

function get_timestamp(date_selector, time_selector) {
  var time = 0
  if (time_selector !== null) {
    time = time_selector.valueAsNumber
  }
  return date_selector.valueAsNumber + time
}

window.onload = _ => {
  /* set up current distance, co2 */
  const current_timeframe_selector = document.getElementById("current_timeframe")
  const current_distance_field = document.getElementById("current_distance")
  const current_co2_field = document.getElementById("current_co2")
  request_value(base_distance_route + "/" + current_timeframe_selector.value, v => current_distance_field.innerHTML = v)
  request_value(base_co2_route + "/" + current_timeframe_selector.value, v => current_co2_field.innerHTML = v)
  current_timeframe_selector.onchange = _ev => {
    var timeframe = current_timeframe_selector.value
    request_value(base_distance_route + "/" + timeframe, v => current_distance_field.innerHTML = v)
    request_value(base_co2_route + "/" + timeframe, v => current_co2_field.innerHTML = v)
  }

  /* set up historic distance, co2 */
  /* limit selectable time range */
  const begin_date_selector = document.getElementById("historic_begin_date")
  const end_date_selector = document.getElementById("historic_end_date")
  const today = new Date()
  var first_historic_day = new Date()
  first_historic_day = new Date(first_historic_day.setDate(today.getDate() - max_days_range))
  begin_date_selector.min = first_historic_day.toISOString().split("T")[0]
  begin_date_selector.max = today.toISOString().split("T")[0]
  end_date_selector.min = first_historic_day.toISOString().split("T")[0]
  end_date_selector.max = today.toISOString().split("T")[0]

  /* Add time selector if necessary */
  const historic_resolution_selector = document.getElementById("historic_resolution")
  if (historic_resolution_selector.value === "1h") {
    add_time_selectors()
  }
  historic_resolution_selector.onchange = _ev => {
    if (historic_resolution_selector.value === "1h") {
      add_time_selectors()
    } else {
      document.getElementById("begin_time_selector").remove()
      document.getElementById("end_time_selector").remove()
    }
  }

  /* manage submit */
  const historic_data_submit_button = document.getElementById("historic_data_submit")
  const historic_distance_field = document.getElementById("historic_distance")
  const historic_co2_field = document.getElementById("historic_co2")
  const used_fields = Array.from(document.querySelectorAll("input[type=date], input[type=time]"))
  historic_data_submit_button.onclick = ev => {
    ev.preventDefault()
    if (used_fields.every(e => e.value !== "")) {
      if (document.getElementById("form_error") !== null) {
        document.getElementById("form_error").remove()
      }
      var params = {
        "begin": get_begin_timestamp(),
        "end": get_end_timestamp(),
        "resolution": historic_resolution_selector.value
      }
      request_value_with_params(base_distance_route + "/history", v => historic_distance_field.value = v, params)
      request_value_with_params(base_co2_route + "/history", v => historic_co2_field.value = v, params)
    } else {
      const error_paragraph = document.createElement("p")
      error_paragraph.id = "form_error"
      error_paragraph.innerHTML = "Missing fields"
      document.getElementById("historic_data_form").append(error_paragraph)
    }
  }
}