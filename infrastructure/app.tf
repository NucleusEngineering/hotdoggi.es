# Disable authentication check to invoke this service
resource "google_cloud_run_service_iam_binding" "app" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.app.name
  role     = "roles/run.invoker"
  members = [
    "allUsers",
  ]
}

# SA with perms for this service
resource "google_service_account" "app" {
  project      = local.project
  account_id   = "${local.prefix}-app"
  display_name = "${local.prefix}-app"
}

# Service definition
resource "google_cloud_run_service" "app" {
  project  = local.project
  provider = google-beta
  name     = "${local.prefix}-app"
  location = local.region
  template {
    spec {
      service_account_name = google_service_account.app.email
      containers {
        image = "gcr.io/${local.project}/app"
        env {
          name  = "GOOGLE_CLOUD_PROJECT"
          value = local.project
        }
      }
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

# Custom domain mapping for this service
resource "google_cloud_run_domain_mapping" "app" {
  project  = local.project
  location = local.region
  name     = "app.${local.domain}"
  metadata {
    namespace = local.project
  }
  spec {
    route_name = google_cloud_run_service.app.name
  }
}

# Service build trigger
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
  description = "Build pipeline for ${local.prefix} app"
  substitutions = {
    _REGION  = local.region
    _PREFIX  = local.prefix
    _SERVICE = "app"
  }

  filename = "services/app/cloudbuild.yaml"
}

# Allow Cloud Build to bind SA
resource "google_service_account_iam_binding" "app-sa-user" {
  service_account_id = google_service_account.app.name
  role               = "roles/iam.serviceAccountUser"
  members = [
    "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
  ]
}