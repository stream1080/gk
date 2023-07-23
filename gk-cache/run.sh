#!/bin/bash
trap "rm gk-cache;kill 0" EXIT

go build -o gk-cache
./gk-cache -port=8001 &
./gk-cache -port=8002 &
./gk-cache -port=8003 -api=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &

wait