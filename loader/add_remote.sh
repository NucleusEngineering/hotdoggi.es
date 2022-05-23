#!/bin/bash

curl -X POST \
    -H "Authorization: Bearer ${TOKEN}" \
    -d @dog.json \
    "api.hotdoggies.stamer.demo.altostrat.com/events/es.hotdoggies.events.dog_added/local-curl"
