#!/bin/sh

curl --include \
     --no-buffer \
     --header "Authorization: Bearer ${TOKEN}" \
     --header "Connection: Upgrade" \
     --header "Upgrade: websocket" \
     --header "Host: api.hotdoggies.stamer.demo.altostrat.com" \
     --header "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" \
     --header "Sec-WebSocket-Version: 13" \
     https://api.hotdoggies.stamer.demo.altostrat.com/dogs/?stream=true