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

resource "google_pubsub_topic" "topic" {
  project = local.project
  name    = "${local.prefix}-stream"
}

resource "google_pubsub_topic" "dead-letter" {
  project = local.project
  name    = "${local.prefix}-dead-letters"
}

resource "google_pubsub_subscription" "subscription" {
  project = local.project
  name    = "${local.prefix}-dead-letter-pull"
  topic   = google_pubsub_topic.dead-letter.name
}

resource "google_project_iam_member" "pubsub-sa-tokencreator" {
  project = local.project
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:service-${local.project_number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "pubsub-sa-publisher" {
  project = local.project
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:service-${local.project_number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "pubsub-sa-subscriber" {
  project = local.project
  role    = "roles/pubsub.subscriber"
  member  = "serviceAccount:service-${local.project_number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_service_account" "pubsub-pusher" {
  project      = local.project
  account_id   = "${local.prefix}-pubsub-pusher"
  display_name = "${local.prefix}-pubsub-pusher"
}