#!/bin/bash

docker run -i owasp/zap2docker-stable zap-baseline.py -t $1/ > web_ui.txt &
docker run -i owasp/zap2docker-stable zap-baseline.py -t $1/airspace/aircraft/list > aircraft_list.txt &
docker run -i owasp/zap2docker-stable zap-baseline.py -t $1/airspace/aircraft/123456/info > aircraft_info.txt &
docker run -i owasp/zap2docker-stable zap-baseline.py -t $1/airspace/history/realtime/1h > realtime_history.txt &
docker run -i owasp/zap2docker-stable zap-baseline.py -t $1/airspace/history > history.txt &
docker run -i owasp/zap2docker-stable zap-baseline.py -t $1/endpoints/position/url > endpoints.txt &
wait
