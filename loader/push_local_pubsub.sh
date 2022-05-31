#!/bin/bash

curl -X POST \
    -d @pubsub.json \
    "http://localhost:8080/events/"
