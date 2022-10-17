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

resource "google_cloud_run_service_iam_member" "public" {
  location = local.region
  project  = local.project
  service  = google_cloud_run_service.proxy.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# SA with perms for this service
resource "google_service_account" "proxy" {
  project      = local.project
  account_id   = "${local.prefix}-proxy"
  display_name = "${local.prefix}-proxy"
}

resource "google_project_iam_member" "proxy-servicecontrol" {
  project = local.project
  role    = "roles/servicemanagement.serviceController"
  member  = "serviceAccount:${google_service_account.proxy.email}"
}

resource "google_project_iam_member" "proxy-sa-cloudtrace" {
  project = local.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.proxy.email}"
}

# Service definition
resource "google_cloud_run_service" "proxy" {
  provider = google-beta
  project  = local.project
  name     = "${local.prefix}-proxy"
  location = local.region
  template {
    spec {
      service_account_name = google_service_account.proxy.email

      timeout_seconds = 3600

      containers {
        image = "gcr.io/${local.project}/proxy"
        resources {
          limits = {
            memory = "512Mi"
            cpu    = "1000m"
          }
        }
      }
    }
  }
  metadata {
    annotations = {
      "run.googleapis.com/ingress"        = "all"
      "client.knative.dev/user-image"     = "gcr.io/${local.project}/proxy"
      "run.googleapis.com/client-name"    = "gcloud"
      "run.googleapis.com/client-version" = local.gcloud_version
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

# Custom domain mapping for this service
resource "google_cloud_run_domain_mapping" "proxy" {
  project  = local.project
  location = local.region
  name     = "api.${local.domain}"
  metadata {
    namespace = local.project
  }
  spec {
    route_name = google_cloud_run_service.proxy.name
  }
}

# Allow Cloud Build to bind SA
resource "google_service_account_iam_member" "proxy-sa-user" {
  provider           = google-beta
  service_account_id = google_service_account.proxy.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
}

resource "google_cloudbuild_trigger" "proxy" {
  project  = local.project
  provider = google-beta
  github {
    name  = local.repo
    owner = local.repo_owner
    push {
      branch = local.branch
    }
  }
  name        = "${local.prefix}-proxy"
  description = "Build pipeline for ${local.prefix}-proxy"
  substitutions = {
    _ENVIRONMENT = "prod"
    _SERVICE     = "proxy"
    _REGION      = local.region
    _PREFIX      = local.prefix
    _API_NAME    = google_endpoints_service.default.service_name
    _API_CONFIG  = google_endpoints_service.default.config_id
  }
  filename = "services/proxy/cloudbuild.yaml"
}

output "gateway_sa" {
  value = google_service_account.proxy.email
}

