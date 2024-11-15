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

resource "google_service_account" "archiver" {
  project      = local.project
  account_id   = "${local.prefix}-archiver"
  display_name = "${local.prefix}-archiver"
}

resource "google_service_account_iam_member" "archiver-sa-user" {
  service_account_id = google_service_account.archiver.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
}

resource "google_project_iam_member" "archiver-sa-cloudtrace" {
  project = local.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.archiver.email}"
}

resource "google_project_iam_member" "archiver-sa-storage" {
  project = local.project
  role    = "roles/storage.objectAdmin"
  member  = "serviceAccount:${google_service_account.archiver.email}"
}

resource "google_cloud_run_service" "archiver" {
  project  = local.project
  provider = google-beta
  name     = "${local.prefix}-archiver"
  location = local.region
  template {
    spec {
      service_account_name = google_service_account.archiver.email
      containers {
        image = "gcr.io/${local.project}/archiver"
        env {
          name  = "ENVIRONMENT"
          value = "prod"
        }
        env {
          name  = "GOOGLE_CLOUD_PROJECT"
          value = local.project
        }
        env {
          name  = "ARCHIVAL_BUCKET"
          value = google_storage_bucket.archiver-bucket.name
        }
        resources {
          limits = {
            memory = "1024Mi"
            cpu    = "1000m"
          }
        }
      }
    }
  }
  metadata {
    annotations = {
      "run.googleapis.com/ingress"        = "all"
      "client.knative.dev/user-image"     = "gcr.io/${local.project}/archiver"
      "run.googleapis.com/client-name"    = "gcloud"
      "run.googleapis.com/client-version" = local.gcloud_version
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "archiver" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.archiver.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.pubsub-pusher.email}"
}

resource "google_pubsub_subscription" "archiver" {
  project = local.project
  name    = "${local.prefix}-archiver-push"
  topic   = google_pubsub_topic.topic.name
  filter  = ""
  push_config {
    push_endpoint = "${google_cloud_run_service.archiver.status[0].url}/v1/events/"
    oidc_token {
      service_account_email = google_service_account.pubsub-pusher.email
    }
  }
  expiration_policy {
    ttl = "99999999s"
  }
  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.dead-letter.id
    max_delivery_attempts = 5
  }
}

resource "google_storage_bucket" "archiver-bucket" {
  project                     = local.project
  name                        = "${local.prefix}-archival-bucket"
  uniform_bucket_level_access = true
  location                    = "EU"
  force_destroy               = true
}

resource "google_cloudbuild_trigger" "archiver" {
  project  = local.project
  provider = google-beta
  github {
    name  = local.repo
    owner = local.repo_owner
    push {
      branch = local.branch
    }
  }
  name        = "${local.prefix}-archiver"
  description = "Build pipeline for ${local.prefix}-archiver"
  substitutions = {
    _ENVIRONMENT     = "prod"
    _SERVICE         = "archiver"
    _REGION          = local.region
    _PREFIX          = local.prefix
    _ARCHIVAL_BUCKET = google_storage_bucket.archiver-bucket.name
  }
  filename = "services/archiver/cloudbuild.yaml"
}
