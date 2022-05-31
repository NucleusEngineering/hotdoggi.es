#!/bin/bash

curl -X POST \
    -H "Authorization: Bearer ${TOKEN}" \
    -d @dog.json \
    "https://api.hotdoggies.stamer.demo.altostrat.com/events/es.hotdoggi.events.dog_moved/local-curl"
