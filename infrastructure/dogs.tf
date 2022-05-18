# Allow Proxy SA to invoke this service
resource "google_cloud_run_service_iam_binding" "dogs" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.dogs.name
  role     = "roles/run.invoker"
  members = [
    "serviceAccount:${google_service_account.proxy.email}"
  ]
}

# SA with perms for this service
resource "google_service_account" "dogs" {
  project      = local.project
  account_id   = "${local.prefix}-dogs"
  display_name = "${local.prefix}-dogs"
}
resource "google_project_iam_member" "dogs-firestore" {
  project = local.project
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.dogs.email}"
}
resource "google_project_iam_member" "dogs-cloudtrace" {
  project = local.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.dogs.email}"
}

# Service definition
resource "google_cloud_run_service" "dogs" {
  project  = local.project
  provider = google-beta
  name     = "${local.prefix}-dogs"
  location = local.region
  template {
    spec {
      service_account_name = google_service_account.dogs.email
      containers {
        image = "gcr.io/${local.project}/dogs"
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
      }
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

# Service build trigger
resource "google_cloudbuild_trigger" "dogs" {
  project  = local.project
  provider = google-beta
  github {
    name  = local.repo
    owner = local.repo_owner
    push {
      branch = local.branch
    }
  }
  name        = "${local.prefix}-dogs"
  description = "Build pipeline for ${local.prefix}-dogs"
  substitutions = {
    _REGION  = local.region
    _PREFIX  = local.prefix
    _SERVICE = "dogs"
  }
  filename = "services/dogs/cloudbuild.yaml"
}

# Allow Cloud Build to bind SA
resource "google_service_account_iam_binding" "dogs-sa-user" {
  service_account_id = google_service_account.dogs.name
  role               = "roles/iam.serviceAccountUser"
  members = [
    "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
  ]
}
