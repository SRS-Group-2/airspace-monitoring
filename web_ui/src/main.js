
const express = require('express')

const app = express()

const port = process.env.PORT || 3000

function range(size, startAt = 0) {
  return [...Array(size).keys()].map(i => i + startAt);
}

app.use(express.static('public'))

app.get('/aircrafts', (_req, res) => {
  res.send(["icao01", "icao02", "icao03", "icao04"])
})

app.get('/aircrafts/:icao24/info', (req, res) => {
  res.send("info about " + req.params.icao24)
})

app.get('/distance/:timeframe((1|6|24)h)', (req, res) => {
  var timeframe = req.params.timeframe
  res.send("distance for last " + timeframe)
})

app.get('/co2/:timeframe((1|6|24)h)', (req, res) => {
  var timeframe = req.params.timeframe
  res.send("co2 for last " + timeframe)
})

app.get('/distance/history', (req, res) => {
  var begin = parseInt(req.query.begin)
  var end = parseInt(req.query.end)
  var resolution = req.query.resolution
  var result = {}
  if (resolution === "1h") {
    var keys1h = range((end - begin) / 3600000, begin / 3600000)
    result = keys1h.reduce((obj, item) => { obj[item] = item; return obj; }, {})
  } else {
    var keys1d = range((end - begin) / 86400000, begin / 86400000)  
    result = keys1d.reduce((obj, item) => { obj[item] = item; return obj; }, {})
  }
  res.json(result)
})

app.get('/co2/history', (req, res) => {
  var begin = parseInt(req.query.begin)
  var end = parseInt(req.query.end)
  var resolution = req.query.resolution
  var result = {}
  if (resolution === "1h") {
    var keys1h = range((end - begin) / 3600000, begin / 3600000)
    result = keys1h.reduce((obj, item) => { obj[item] = item; return obj; }, {})
  } else {
    var keys1d = range((end - begin) / 86400000, begin / 86400000)  
    result = keys1d.reduce((obj, item) => { obj[item] = item; return obj; }, {})
  }
  res.json(result)
})

app.get("/test", (_req, res) => {
  res.send("Happy use of our service!")
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})

