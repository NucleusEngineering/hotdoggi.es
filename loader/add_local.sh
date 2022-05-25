#!/bin/bash

curl -X POST \
    -d @dog.json \
    "http://localhost:8080/events/es.hotdoggi.events.dog_added/local-curl"
