const base_aircraft_route = "/airspace/aircraft"
const base_history_route = "/airspace/history"

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

function set_realtime_history(timeframe) {
  request_value(base_history_route + "/history/realtime/" + timeframe, v => {
    document.getElementById("current_distance").innerHTML = v // TODO:
  })
}

function show_history(history) {
  document.getElementById("historic_distance").innerHTML = history // TODO:
}

function set_selectable_italian_flights() {
  const flight_selector = document.getElementById("flight")
  request_value(base_aircraft_route + "/list", vs => {
    flight_selector.innerHTML = ""
    JSON.parse(vs).map(v => {
      var opt = document.createElement("option")
      opt.value = v
      opt.innerHTML = v
      flight_selector.append(opt)
    })
  })
}

window.onload = _ => {
  /* set up flight selector */
  const flight_selector = document.getElementById("flight")
  set_selectable_italian_flights()
  flight_selector.onclick = _ev => {
    set_selectable_italian_flights()
  }
  flight_selector.onchange = _ev => {
    request_value(base_aircraft_route + "/" + flight_selector.value + "/info", v => {
      document.getElementById("info").innerHTML = v
      document.getElementById("coordinates").innerHTML = "Coordinates of " + flight_selector.value
      // TODO with websockets
    })
  }
  /* set up current distance, co2 */
  const current_timeframe_selector = document.getElementById("current_timeframe")
  set_realtime_history(current_timeframe_selector.value)
  current_timeframe_selector.onchange = _ev => {
    set_realtime_history(current_timeframe_selector.value)
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
  historic_data_submit_button.onclick = ev => {
    ev.preventDefault()
    const used_fields = Array.from(document.querySelectorAll("input[type=date], input[type=time]"))
    if (used_fields.every(e => e.value !== "")) {
      if (document.getElementById("form_error") !== null) {
        document.getElementById("form_error").remove()
      }
      var params = {
        "from": get_begin_timestamp(),
        "to": get_end_timestamp(),
        "resolution": historic_resolution_selector.value
      }
      request_value_with_params(base_history_route + "/history", show_history, params)
    } else {
      const error_paragraph = document.createElement("p")
      error_paragraph.id = "form_error"
      error_paragraph.innerHTML = "Missing fields"
      document.getElementById("historic_data_form").append(error_paragraph)
    }
  }
}