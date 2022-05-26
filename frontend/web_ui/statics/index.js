const base_aircraft_route = "/airspace/aircraft"
const base_history_route = "/airspace/history"

const max_days_range = 30

function request_value(url, value_callback) {
  var xhttp = new XMLHttpRequest()
  xhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      value_callback(xhttp.responseText)
    }
    else if (this.readyState == 4 && this.status == 404)
    {
      value_callback("not found")
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

  return (date_selector.valueAsNumber + time)/1000
}

function set_realtime_history(timeframe) {
  
  document.getElementById("current_error").innerText=""
  
  request_value(base_history_route + "/realtime/" + timeframe, v => {
    if(v=="{}") {
      document.getElementById("current_error").innerText = "No data available"
      document.getElementById("km").innerText=""
      document.getElementById("co2").innerText=""
      document.getElementById("current_update").innerText=""
    } else {
      var json = JSON.parse(v)
      const current_distance = document.getElementById("km")
      const current_co2 = document.getElementById("co2")
      document.getElementById("current_update").innerHTML =  get_date_string(json.timestamp)
      current_distance.innerText = json.distanceKm
      current_co2.innerText = json.CO2t
    }
  })
}

function show_daily_history(history) {
  var json = JSON.parse(history)
  table = document.getElementById("historic_distance")
  table.innerHTML=""
  table= document.createElement("table")
  table.setAttribute('id','historical-data')
  var h_row = document.createElement("tr")
  var h_time = document.createElement("th")
  h_time.innerHTML = "Time"
  var h_distance = document.createElement("th")
  h_distance.innerHTML = "Distance (km)"
  var h_co2 = document.createElement("th")
  h_co2.innerHTML = "CO2 (t)"
  h_row.append(h_time, h_distance, h_co2)
  var rowDel = document.createElement("tr")
  rowDel.setAttribute("id","toDelete")
  for(i=0;i<=2;i++) {
    var toAdd=document.createElement("td")
    rowDel.append(toAdd)
  }
  table.append(h_row,rowDel)
  document.getElementById("historic_distance").append(table)
  if(Object.entries(json.history).length===0) {
    const error_paragraph = document.createElement("p")
    error_paragraph.id = "form_error"
    error_paragraph.innerHTML = "No data available in this interval"
    error_paragraph.setAttribute("class","text")
    document.getElementById("historic_data_form").append(error_paragraph)
  } else {
    Object.entries(json.history).forEach(element => {
      document.getElementById("toDelete").innerHTML=""
      var time = document.createElement("td")
      var row = document.createElement("tr")
      if(json.resolution=="day") {
        time.innerHTML= get_onlydate_string(element[1]["startTime"])
      }
      else if (json.resolution=="hour") {
        time.innerHTML = get_date_string(element[1]["startTime"])
      }
      var distance = document.createElement("td")
      distance.innerHTML = element[1]["distanceKm"]
      var co2 = document.createElement("td")
      co2.innerHTML = element[1]["CO2t"]
      row.append(time, distance, co2)
      table.append(row)
    })
    document.getElementById("historic_distance").append(table)
  }
}



function get_date_string(timestamp) {
  var args = timestamp.split("-").map(s => parseInt(s))
  args[1]=args[1]-1
  return (new Date(Date.UTC(...args))).toLocaleString()
}

function get_onlydate_string(timestamp) {
  var args = timestamp.split("-").map(s => parseInt(s))
  args[1]=args[1]-1
  return (new Date(Date.UTC(...args))).toLocaleDateString()
}

function set_selectable_italian_flights() {
  const flight_selector = document.getElementById("flight")
  if(flight_selector.value.length==6){
    flight=flight_selector.value
  }
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
      opt.setAttribute("id",v)
      opt.value = v
      opt.innerHTML = v
      flight_selector.append(opt)
    })
    if(flight.length==6) {
      document.getElementById(flight).setAttribute("selected",true)
    }
  })
}

function set_aircraft_info(info) {
  if(info=="not found") {
    document.getElementById("notFound").innerText="Information about this icao24 was not found in our database"
    document.getElementById("icao").innerText=""
    document.getElementById("manufacturer").innerText=""
    document.getElementById("model").innerText = ""
    document.getElementById("registration-n").innerText=""
    document.getElementById("serial-number").innerText = ""
  }
  var json = JSON.parse(info)
  document.getElementById("icao").innerText=json.icao24
  document.getElementById("manufacturer").innerText=json.manufacturer
  document.getElementById("model").innerText = json.model
  document.getElementById("registration-n").innerText=json.registration
  document.getElementById("serial-number").innerText = json.serialnumber  
}

function update_aircraft_position(data) {
  var json = JSON.parse(data)
  var time = new Date(json.timestamp * 1000).toLocaleString()
  document.getElementById("icao24").innerText=json.icao24
  document.getElementById("latitude").innerText = json.lat
  document.getElementById("longitude").innerText=json.lon
  document.getElementById("updatetime").innerText = time  
}

window.onload = _ => {
  console.log("Loaded")
  document.getElementById("loading").innerText=""
  var websocket = null
  /* set up flight selector */
  const flight_selector = document.getElementById("flight")
  flight_selector.onclick = _ev => {
    set_selectable_italian_flights()
  }
  flight_selector.onchange = _ev => {
    document.getElementById("startedListening").innerText="Started listening positions at "+ new Date(Date.now()).toLocaleString()
    document.getElementById("notFound").innerText=""
    document.getElementById("icao24").innerText=""
    document.getElementById("latitude").innerText=""
    document.getElementById("longitude").innerText=""
    document.getElementById("updatetime").innerText=""
    if (websocket !== null) {
        websocket.close()
    }
    flight=flight_selector.value
    request_value(base_aircraft_route + "/" + flight + "/info", v => {
      set_aircraft_info(v)})
      request_value("/endpoints/position/url", v => {
      var url = v.replace("https://", "wss://") + base_aircraft_route + "/" + flight + "/position"
      websocket = new WebSocket(url)
      websocket.onopen = function(event) {
        document.getElementById("loading").innerText="Waiting for next position..."
      }
      websocket.addEventListener('message', ev =>  {
          update_aircraft_position(ev.data)
      })
      websocket.onclose = function (event) {
          document.getElementById("loading").innerText="No more position available"
      }
    })
  }
  /* set up current distance, co2 */
  const current_timeframe_selector = document.getElementById("current_timeframe")
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
      request_value_with_params(base_history_route, show_daily_history, params)
    } else {
      const error_paragraph = document.createElement("p")
      error_paragraph.id = "form_error"
      error_paragraph.innerHTML = "Missing fields"
      error_paragraph.setAttribute("class","text")
      document.getElementById("historic_data_form").append(error_paragraph)
    }
  }
}
