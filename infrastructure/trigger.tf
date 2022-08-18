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

resource "google_service_account" "trigger" {
  project      = local.project
  account_id   = "${local.prefix}-trigger"
  display_name = "${local.prefix}-trigger"
}

resource "google_service_account_iam_member" "trigger-sa-user" {
  service_account_id = google_service_account.trigger.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
}

resource "google_project_iam_member" "trigger-sa-logging" {
  project = local.project
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.trigger.email}"
}

resource "google_project_iam_member" "trigger-sa-cloudtrace" {
  project = local.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.trigger.email}"
}

resource "google_project_iam_member" "trigger-sa-firestore" {
  project = local.project
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.trigger.email}"
}

resource "google_project_iam_member" "trigger-sa-pubsub" {
  project = local.project
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${google_service_account.trigger.email}"
}

resource "google_project_iam_member" "trigger-sa-pubsub-view" {
  project = local.project
  role    = "roles/pubsub.viewer"
  member  = "serviceAccount:${google_service_account.trigger.email}"
}

resource "google_cloudfunctions_function" "function" {
  project               = local.project
  name                  = "${local.prefix}-trigger"
  description           = "${local.prefix}-trigger"
  runtime               = "ruby27"
  service_account_email = google_service_account.trigger.email
  available_memory_mb   = 256

  ingress_settings = "ALLOW_INTERNAL_AND_GCLB"

  source_archive_bucket = google_storage_bucket.function-bucket.name
  source_archive_object = google_storage_bucket_object.function-source.name

  entry_point = "function"
  event_trigger {
    event_type = "providers/cloud.firestore/eventTypes/document.create"
    resource   = "projects/${local.project}/databases/(default)/documents/{collection}/{event}"
  }

  environment_variables = {
    TOPIC                = google_pubsub_topic.topic.name,
    GOOGLE_CLOUD_PROJECT = local.project
  }
}

resource "google_cloudbuild_trigger" "trigger" {
  project  = local.project
  provider = google-beta
  github {
    name  = local.repo
    owner = local.repo_owner
    push {
      branch = local.branch
    }
  }
  name        = "${local.prefix}-trigger"
  description = "Build pipeline for ${local.prefix}-trigger"

  substitutions = {
    _ENVIRONMENT = "prod"
    _FUNCTION    = "trigger"
    _PREFIX      = local.prefix
    _SOURCE      = "gs://${google_storage_bucket.function-bucket.name}/${google_storage_bucket_object.function-source.name}"
    _REGION      = local.region
    _BUCKET      = google_storage_bucket.function-bucket.name
  }

  filename = "services/trigger/cloudbuild.yaml"
}

resource "google_project_iam_member" "trigger-function-deployer" {
  project = local.project
  role    = "roles/cloudfunctions.admin"
  member  = "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
}

resource "google_storage_bucket" "function-bucket" {
  project                     = local.project
  name                        = "${local.prefix}-function-source-bucket"
  uniform_bucket_level_access = true
  location                    = "EU"
  force_destroy               = true
}

resource "google_storage_bucket_object" "function-source" {
  name   = "function.zip"
  bucket = google_storage_bucket.function-bucket.name
  source = "../services/trigger/function.zip"
}

output "function_bucket" {
  value = google_storage_bucket.function-bucket.name
}