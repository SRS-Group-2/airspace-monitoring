
const express = require('express')
const app = express()
const port = 3000

function range(size, startAt = 0) {
  return [...Array(size).keys()].map(i => i + startAt);
}

app.use(express.static('public'))

app.get('/distance/:timeframe((1|6|24)h)', (req, res) => {
  var timeframe = req.params.timeframe
  res.send("distance for last " + timeframe)
})

app.get('/co2/:timeframe((1|6|24)h)', (req, res) => {
  var timeframe = req.params.timeframe
  res.send("co2 for last " + timeframe)
})

app.get('/distance/history', (req, res) => {
  console.log(req)
  var begin = req.params.begin
  var end = req.params.end
  var resolution = req.params.resolution
  var result = {}
  if (resolution === "1h") {
    var keys1h = range((end - begin) / 3600000, begin / 3600000)  
    result = keys1h.reduce((obj, item) => { return {...obj, item: item} })
  } else {
    var keys1d = range((end - begin) / 86400000, begin / 86400000)  
    result = keys1d.reduce((obj, item) => { return {...obj, item: item} })
  }
  res.send(result)
})

app.get('/co2/history', (req, res) => {
  var begin = req.params.begin
  var end = req.params.end
  var resolution = req.params.resolution
  var result = {}
  if (resolution === "1h") {
    var keys1h = range((end - begin) / 3600000, begin / 3600000)  
    result = keys1h.reduce((obj, item) => { return {...obj, item: item} })
  } else {
    var keys1d = range((end - begin) / 86400000, begin / 86400000)  
    result = keys1d.reduce((obj, item) => { return {...obj, item: item} })
  }
  res.send(result)
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})

