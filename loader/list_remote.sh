#!/bin/bash

curl -X GET \
    -H "Authorization: Bearer ${TOKEN}" \
    "api.hotdoggies.stamer.demo.altostrat.com/dogs/"
