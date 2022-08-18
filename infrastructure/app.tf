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

resource "google_storage_bucket" "app" {
  project                     = local.project
  provider                    = google-beta
  name                        = "${local.prefix}-app-bucket"
  uniform_bucket_level_access = true
  location                    = "EU"

  force_destroy = true
  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
}

resource "google_compute_backend_bucket" "app" {
  project     = local.project
  provider    = google-beta
  name        = "${local.prefix}-app-backend"
  description = "Origin for static SPA serving"
  bucket_name = google_storage_bucket.app.name
  enable_cdn  = true
}

resource "google_storage_bucket_iam_member" "app" {
  bucket = google_storage_bucket.app.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

resource "google_cloudbuild_trigger" "app" {
  project  = local.project
  provider = google-beta
  github {
    name  = local.repo
    owner = local.repo_owner
    push {
      branch = local.branch
    }
  }
  name        = "${local.prefix}-app"
  description = "Build pipeline for ${local.prefix}-app"
  substitutions = {
    _REGION  = local.region
    _PREFIX  = local.prefix
    _SERVICE = "app"
    _BUCKET  = google_storage_bucket.app.name
  }

  filename = "app/cloudbuild.yaml"
}
