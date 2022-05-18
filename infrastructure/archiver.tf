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

resource "google_project_iam_member" "archiver-sa-logging" {
  project = local.project
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.archiver.email}"
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

resource "google_pubsub_subscription" "archiver" {
  project = local.project
  name    = "${local.prefix}-archiver-push"
  topic   = google_pubsub_topic.topic.name
  filter  = ""
  push_config {
    push_endpoint = google_cloud_run_service.archiver.status[0].url
    oidc_token {
      service_account_email = google_service_account.pubsub-pusher.email
    }
  }
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
          name  = "GOOGLE_CLOUD_PROJECT"
          value = local.project
        }
        env {
          name  = "HOTDOGGIES_ARCHIVAL_BUCKET"
          value = google_storage_bucket.archiver-bucket.name
        }
        resources {
          limits = {
            memory = "1024Mi"
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

resource "google_cloud_run_service_iam_binding" "archiver-users" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.archiver.name
  role     = "roles/run.invoker"
  members = [
    "serviceAccount:${google_service_account.pubsub-pusher.email}"
  ]
}

resource "google_storage_bucket" "archiver-bucket" {
  project                     = local.project
  name                        = "${local.prefix}-archive-bucket"
  uniform_bucket_level_access = true
  location                    = "EU"
  force_destroy               = true
}

resource "google_cloudbuild_trigger" "archiver" {
  project     = local.project
  provider    = google-beta
  name        = "${local.prefix}-archiver"
  description = "${local.prefix}-archiver-ci"
  substitutions = {
    _HOTDOGGIES_ENVIRONMENT     = "prod"
    _HOTDOGGIES_SERVICE         = "archiver"
    _HOTDOGGIES_REGION          = local.region
    _HOTDOGGIES_PREFIX          = local.prefix
    _HOTDOGGIES_ARCHIVAL_BUCKET = google_storage_bucket.archiver-bucket.name
  }
  filename = "../services/archiver/cloudbuild.yaml"

  trigger_template {
    project_id  = local.project
    branch_name = "main"
    repo_name   = local.repo
  }
}

output "archiver-endpoint" {
  value = google_cloud_run_service.archiver.status[0].url
}
