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

resource "google_project_iam_member" "ingest-sa-logging" {
  project = local.project
  role    = "roles/logging.logWriter"
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
  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_binding" "ingest-users" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.ingest.name
  role     = "roles/run.invoker"
  members = [
    "user:stamer@google.com"
  ]
}

resource "google_cloudbuild_trigger" "ingest" {
  project     = local.project
  provider    = google-beta
  name        = "${local.prefix}-ingest"
  description = "${local.prefix}-ingest-ci"
  substitutions = {
    _HOTDOGGIES_ENVIRONMENT = "prod"
    _HOTDOGGIES_SERVICE     = "ingest"
    _HOTDOGGIES_REGION      = local.region
    _HOTDOGGIES_PREFIX      = local.prefix
  }
  filename = "../services/ingest/cloudbuild.yaml"

  trigger_template {
    project_id  = local.project
    branch_name = "main"
    repo_name   = local.repo
  }
}

output "ingest-endpoint" {
  value = google_cloud_run_service.ingest.status[0].url
}
