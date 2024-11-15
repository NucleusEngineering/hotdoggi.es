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

resource "google_service_account" "analytics" {
  project      = local.project
  account_id   = "${local.prefix}-analytics"
  display_name = "${local.prefix}-analytics"
}

resource "google_service_account_iam_member" "analytics-sa-user" {
  service_account_id = google_service_account.analytics.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
}

resource "google_project_iam_member" "analytics-sa-cloudtrace" {
  project = local.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.analytics.email}"
}

resource "google_project_iam_member" "analytics-sa-bq" {
  project = local.project
  role    = "roles/bigquery.jobUser"
  member  = "serviceAccount:${google_service_account.analytics.email}"
}

resource "google_cloud_run_service" "analytics" {
  project  = local.project
  provider = google-beta
  name     = "${local.prefix}-analytics"
  location = local.region
  template {
    spec {
      service_account_name = google_service_account.analytics.email
      containers {
        image = "gcr.io/${local.project}/analytics"
        env {
          name  = "ENVIRONMENT"
          value = "prod"
        }
        env {
          name  = "GOOGLE_CLOUD_PROJECT"
          value = local.project
        }
        env {
          name  = "BQ_DATASET"
          value = google_bigquery_dataset.dataset.dataset_id
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
      "client.knative.dev/user-image"     = "gcr.io/${local.project}/analytics"
      "run.googleapis.com/client-name"    = "gcloud"
      "run.googleapis.com/client-version" = local.gcloud_version
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "analytics" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.analytics.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.pubsub-pusher.email}"
}

resource "google_pubsub_subscription" "analytics" {
  project = local.project
  name    = "${local.prefix}-analytics-push"
  topic   = google_pubsub_topic.topic.name
  filter  = ""
  push_config {
    push_endpoint = "${google_cloud_run_service.analytics.status[0].url}/v1/events/"
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

resource "google_cloudbuild_trigger" "analytics" {
  project  = local.project
  provider = google-beta
  github {
    name  = local.repo
    owner = local.repo_owner
    push {
      branch = local.branch
    }
  }
  name        = "${local.prefix}-analytics"
  description = "Build pipeline for  ${local.prefix}-analytics"
  substitutions = {
    _ENVIRONMENT = "prod"
    _SERVICE     = "analytics"
    _REGION      = local.region
    _PREFIX      = local.prefix
    _BQ_DATASET  = google_bigquery_dataset.dataset.dataset_id
  }
  filename = "services/analytics/cloudbuild.yaml"
}

resource "google_bigquery_dataset" "dataset" {
  project       = local.project
  dataset_id    = "${local.prefix}_events"
  friendly_name = "${local.prefix} events"
  description   = "All received events"
  location      = "EU"

  labels = {
    env = "default"
  }

  access {
    role          = "OWNER"
    user_by_email = google_service_account.analytics.email
  }

  access {
    role          = "OWNER"
    user_by_email = local.user
  }
}

