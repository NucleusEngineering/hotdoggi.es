# Copyright 2022 Google

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/sh

service_name="$(terraform -chdir=../../infrastructure output -json | jq -r '.service_name.value')"
config_id="$(terraform -chdir=../../infrastructure output -json | jq -r '.config_id.value')"

token="$(gcloud auth print-access-token)"

service_config=$(curl -H "Authorization: Bearer ${token}" \
  "https://servicemanagement.googleapis.com/v1/services/${service_name}/configs/${config_id}?view=FULL")

echo $service_config > service.json
