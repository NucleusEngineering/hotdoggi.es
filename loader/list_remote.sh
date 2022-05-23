#!/bin/bash

curl -X GET \
    -H "Authorization: Bearer ${TOKEN}" \
    "https://api.hotdoggies.stamer.demo.altostrat.com/dogs/"
