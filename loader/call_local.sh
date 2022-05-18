#!/bin/bash

curl -X POST \
    -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
    -d @payload.json \
    "localhost:8080/com.google.corp.events.shipment.delivered/python-loader"
