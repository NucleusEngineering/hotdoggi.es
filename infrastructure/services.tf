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
resource "google_project_service" "endpoints" {
  project            = local.project
  service            = "endpoints.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "servicecontrol" {
  project            = local.project
  service            = "servicecontrol.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "servicemanagement" {
  project            = local.project
  service            = "servicemanagement.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "compute" {
  project            = local.project
  service            = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "cloudrun" {
  project            = local.project
  service            = "run.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "cloudfunctions" {
  project            = local.project
  service            = "cloudfunctions.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "cloudbuild" {
  project            = local.project
  service            = "cloudbuild.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "cloudtrace" {
  project            = local.project
  service            = "cloudtrace.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "firestore" {
  project            = local.project
  service            = "firestore.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "identitytoolkit" {
  project            = local.project
  service            = "identitytoolkit.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "securetoken" {
  project            = local.project
  service            = "securetoken.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "dns" {
  project            = local.project
  service            = "dns.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "domains" {
  project            = local.project
  service            = "domains.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "appengine" {
  project            = local.project
  service            = "appengine.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "cloudresourcemanager" {
  project            = local.project
  service            = "cloudresourcemanager.googleapis.com"
  disable_on_destroy = false
}

# Allow Cloud Build to deploy to Cloud Run
resource "google_project_iam_member" "cloudbuild-deploy" {
  project = local.project
  role    = "roles/run.admin"
  member  = "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
}
# Allow Cloud Build to read the Service Management API
resource "google_project_iam_member" "cloudbuild-servicecontroller" {
  project = local.project
  role    = "roles/servicemanagement.serviceController"
  member  = "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
}