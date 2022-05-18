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
  request_value(base_history_route + "/realtime/" + timeframe, v => {
    console.log(v)
    var json = JSON.parse(v) // strange error here
    const current_distance = document.getElementById("current_distance")
    const current_co2 = document.getElementById("current_co2")
    document.getElementById("current_update") = "Data last updated at " + get_date_string(json.timestamp)
    current_distance.innerHTML = "Distance travelled in the last " + timeframe + " is " + json["distanceKm"]
    current_co2.innerHTML = "CO2 emitted in the last " + timeframe + " is " + json["CO2t"]
  })
}

function show_history(history) {
  var json = JSON.parse(history)
  var table = document.createElement(table)
  var h_row = document.createElement("tr")
  var h_time = document.createElement("td")
  h_time.innerHTML = "Time"
  var h_distance = document.createElement("td")
  h_distance.innerHTML = "Distance (km)"
  var h_co2 = document.createElement("td")
  h_co2.innerHTML = "CO2 (t)"
  h_row.append(h_time, h_distance, h_co2)
  table.append(h_row)
  Object.entries(json).forEach(element => {
    var row = document.createElement("tr")
    var time = document.createElement("td")
    time.innerHTML = get_date_string(element[1]["startTime"])
    var distance = document.createElement("td")
    distance.innerHTML = element[1]["distanceKm"]
    var co2 = document.createElement("td")
    co2.innerHTML = element[1]["CO2t"]
    row.append(time, distance, co2)
    table.append(row)
  })
  document.getElementById("historic_distance").innerHTML = ""
  document.getElementById("historic_distance").append(table)
}

function get_date_string(timestamp) {
  var args = timestamp.split("-").map(s => parseInt(s))
  return (new Date(...args)).toLocaleString()
}

function set_selectable_italian_flights() {
  const flight_selector = document.getElementById("flight")
  request_value(base_aircraft_route + "/list", vs => {
    var json = JSON.parse(vs)
    document.getElementById("flights_list_update").innerHTML = "Aircraft list updated at " + get_date_string(json.timestamp)
    flight_selector.innerHTML = ""
    var emptyOption = document.createElement("option")
    emptyOption.value = " "
    emptyOption.innerHTML = " "
    flight_selector.append(emptyOption)
    json["icao24"].map(v => {
      var opt = document.createElement("option")
      opt.value = v
      opt.innerHTML = v
      flight_selector.append(opt)
    })
  })
}

function set_aircraft_info(info) {
  var json = JSON.parse(info)
  document.getElementById("info").innerHTML = "Manufacturer: " + json.manufacturer + "<br>"
                                            + "Model: " + json.model + "<br>"
                                            + "Registration number: " + json.registration + "<br>"
                                            + "Serial number: " + json.serialnumber
}

function update_aircraft_position(data) {
  var json = JSON.parse(data)
  var time = new Date(json.timestamp * 1000).toLocaleString()
  document.getElementById("coordinates_data").innerHTML = "Position updated at " + time + "<br>"
                                                        + "Latitude: " + json.lat + "<br>"
                                                        + "Longitude: " + json.lon
}

window.onload = _ => {
  var websocket = null
  /* set up flight selector */
  const flight_selector = document.getElementById("flight")
  set_selectable_italian_flights()
  flight_selector.onclick = _ev => {
    set_selectable_italian_flights()
  }
  flight_selector.onchange = _ev => {
    document.getElementById("coordinates_data").innerHTML =""
    
    if (flight_selector.value === " ") {
      document.getElementById("info").innerHTML = ""
      if (websocket !== null) {
        websocket.close()
      }
    } else {
      flight=flight_selector.value
      request_value(base_aircraft_route + "/" + flight + "/info", v => {
        set_aircraft_info(v)})
      document.getElementById("coordinates_title").innerHTML = "Current position of aircraft " + flight
      // TODO: improve websocket and websocket error management
      var url = window.location.href.slice(0, -1).replace("https://", "ws://") + base_aircraft_route + "/" + flight + "/position"
      websocket = new WebSocket(url)
      websocket.addEventListener('message', ev =>  {
      update_aircraft_position(ev.data)
        })
      websocket.onclose = function (event) {
        if(event.code == 1006){
          alert("The connection was closed abnormally: code 1006")
        }
      }
    }
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
  if (historic_resolution_selector.value === "hour") {
    add_time_selectors()
  }
  historic_resolution_selector.onchange = _ev => {
    if (historic_resolution_selector.value === "hour") {
      add_time_selectors()
    } else {
      document.getElementById("begin_time_selector").remove()
      document.getElementById("end_time_selector").remove()
    }
  }

  /* manage submit */
  const historic_data_submit_button = document.getElementById("historic_data_submit")
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