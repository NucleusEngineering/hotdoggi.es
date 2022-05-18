#!/bin/sh

cd ../../infrastructure/
service_name="$(terraform output -json | jq -r '.service_name.value')"
config_id="$(terraform output -json | jq -r '.config_id.value')"

token="$(gcloud auth print-access-token)"

service_config=$(curl -H "Authorization: Bearer ${token}" \
  "https://servicemanagement.googleapis.com/v1/services/${service_name}/configs/${config_id}?view=FULL")

echo $service_config > ../services/proxy/service.json
