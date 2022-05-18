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

  source_repository {
    url = "https://source.developers.google.com/projects/${local.project}/repos/${local.repo}/moveable-aliases/main/paths/trigger"
  }

  entry_point = "function"
  event_trigger {
    event_type = "providers/cloud.firestore/eventTypes/document.create"
    resource   = "projects/${local.project}/databases/(default)/documents/{collection}/{event}"
  }

  environment_variables = {
    HOTDOGGIES_TOPIC           = google_pubsub_topic.topic.name,
    GOOGLE_CLOUD_PROJECT = local.project
  }
}

resource "google_cloudbuild_trigger" "trigger" {
  project     = local.project
  provider    = google-beta
  name        = "${local.prefix}-trigger"
  description = "${local.prefix}-trigger-ci"

  substitutions = {
    _HOTDOGGIES_ENVIRONMENT = "prod"
    _HOTDOGGIES_FUNCTION    = "trigger"
    _HOTDOGGIES_PREFIX      = local.prefix
    _HOTDOGGIES_REPO        = google_cloudfunctions_function.function.source_repository[0].url
    _HOTDOGGIES_REGION      = local.region
  }

  filename = "../services/trigger/cloudbuild.yaml"

  trigger_template {
    project_id  = local.project
    branch_name = "main"
    repo_name   = local.repo
  }
}
