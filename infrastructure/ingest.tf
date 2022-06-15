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

resource "google_service_account" "ingest" {
  project      = local.project
  account_id   = "${local.prefix}-ingest"
  display_name = "${local.prefix}-ingest"
}

resource "google_service_account_iam_member" "ingest-sa-user" {
  service_account_id = google_service_account.ingest.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
}

resource "google_project_iam_member" "ingest-sa-firestore" {
  project = local.project
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.ingest.email}"
}

resource "google_project_iam_member" "ingest-sa-cloudtrace" {
  project = local.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.ingest.email}"
}

resource "google_cloud_run_service" "ingest" {
  project  = local.project
  provider = google-beta
  name     = "${local.prefix}-ingest"
  location = local.region
  template {
    spec {
      service_account_name = google_service_account.ingest.email
      containers {
        image = "gcr.io/${local.project}/ingest"
        env {
          name  = "GATEWAY_SA"
          value = google_service_account.proxy.email
        }
        env {
          name  = "ENVIRONMENT"
          value = "prod"
        }
        env {
          name  = "GOOGLE_CLOUD_PROJECT"
          value = local.project
        }
        resources {
          limits = {
            memory = "256Mi"
          }
        }
      }
    }
  }
  metadata {
    annotations = {
      "run.googleapis.com/ingress" = "all"
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "ingest" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.ingest.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.proxy.email}"
}

resource "google_cloudbuild_trigger" "ingest" {
  project  = local.project
  provider = google-beta
  github {
    name  = local.repo
    owner = local.repo_owner
    push {
      branch = local.branch
    }
  }
  name        = "${local.prefix}-ingest"
  description = "Build pipeline for ${local.prefix}-ingest"
  substitutions = {
    _ENVIRONMENT = "prod"
    _SERVICE     = "ingest"
    _REGION      = local.region
    _PREFIX      = local.prefix
  }
  filename = "services/ingest/cloudbuild.yaml"
}
