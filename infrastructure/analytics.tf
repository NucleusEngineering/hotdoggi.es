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

resource "google_project_iam_member" "analytics-sa-logging" {
  project = local.project
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.analytics.email}"
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

resource "google_bigquery_dataset" "dataset" {
  project                     = local.project
  dataset_id                  = "${local.prefix}_events"
  friendly_name               = "${local.prefix} events"
  description                 = "All received HOTDOGGIES events"
  location                    = "EU"

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

  access {
    role   = "READER"
    domain = local.domain
  }
}

resource "google_pubsub_subscription" "analytics" {
  project = local.project
  name    = "${local.prefix}-analytics-push"
  topic   = google_pubsub_topic.topic.name
  filter  = ""
  push_config {
    push_endpoint = google_cloud_run_service.analytics.status[0].url
    oidc_token {
      service_account_email = google_service_account.pubsub-pusher.email
    }
  }
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
          name  = "GOOGLE_CLOUD_PROJECT"
          value = local.project
        }
        env {
          name  = "HOTDOGGIES_BQ_DATASET"
          value = google_bigquery_dataset.dataset.dataset_id
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

resource "google_cloud_run_service_iam_binding" "analytics-users" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.analytics.name
  role     = "roles/run.invoker"
  members = [
    "serviceAccount:${google_service_account.pubsub-pusher.email}"
  ]
}

resource "google_cloudbuild_trigger" "analytics" {
  project     = local.project
  provider    = google-beta
  name        = "${local.prefix}-analytics"
  description = "${local.prefix}-analytics-ci"
  substitutions = {
    _HOTDOGGIES_ENVIRONMENT     = "prod"
    _HOTDOGGIES_SERVICE         = "analytics"
    _HOTDOGGIES_REGION          = local.region
    _HOTDOGGIES_PREFIX          = local.prefix
    _HOTDOGGIES_BQ_DATASET      = google_bigquery_dataset.dataset.dataset_id
  }
  filename = "../services/analytics/cloudbuild.yaml"

  trigger_template {
    project_id  = local.project
    branch_name = "main"
    repo_name   = local.repo
  }
}

output "analytics-endpoint" {
  value = google_cloud_run_service.analytics.status[0].url
}
