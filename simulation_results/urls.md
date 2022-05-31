### Base URL:
https://airspace-monitoring-5xotb2mk.ew.gateway.dev

### WS URL:
wss://aircraft-positions-onwhlgspyq-ew.a.run.app

### Endpoints:
- /index.html
- /airspace.css
- /index.js
- /favicon.ico
- /hero-bg.png
- /plane.png
- /airspace/aircraft/list
- /airspace/aircraft/123456/info
- /airspace/history/realtime/6h
- /airspace/history?from=1653733500&to=1653992700&resolution=hour
- /endpoints/position/url
- ws://aircraft-positions-onwhlgspyq-ew.a.run.app/airspace/aircraft/:icao24/position

### Test format:
hey.exe -c 100 -z 60s https://airspace-monitoring-5xotb2mk.ew.gateway.dev/airspace/aircraft/list