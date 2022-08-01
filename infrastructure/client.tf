# Copyright 2022 Daniel Stamer

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

resource "google_storage_bucket" "client" {
  project                     = local.project
  name                        = "${local.prefix}-client-artifacts-bucket"
  uniform_bucket_level_access = true
  location                    = "EU"
  force_destroy               = true
}

resource "google_cloudbuild_trigger" "client" {
  project  = local.project
  provider = google-beta
  github {
    name  = local.repo
    owner = local.repo_owner
    push {
      branch = local.branch
    }
  }
  name        = "${local.prefix}-client"
  description = "Build pipeline for the websocket client"
  substitutions = {
    _BUCKET = google_storage_bucket.client.name
  }
  filename = "clients/dogs/cloudbuild.yaml"
}
