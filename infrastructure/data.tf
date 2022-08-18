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

# Firestore database for this project
resource "google_app_engine_application" "database" {
  project       = local.project
  provider      = google-beta
  location_id   = replace(replace(replace(replace(local.region, "1", ""), "2", ""), "3", ""), "4", "") # TODO strip all trailing numbers
  database_type = "CLOUD_FIRESTORE"
}
