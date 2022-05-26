#!/bin/bash

curl -X POST \
    -d @event.json \
    "http://localhost:8080/events/"
