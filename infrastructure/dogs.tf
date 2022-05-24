resource "google_service_account" "dogs" {
  project      = local.project
  account_id   = "${local.prefix}-dogs"
  display_name = "${local.prefix}-dogs"
}

resource "google_service_account_iam_member" "dogs-sa-user" {
  service_account_id = google_service_account.dogs.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
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

resource "google_cloud_run_service_iam_member" "dogs-gateway" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.dogs.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.proxy.email}"
}

resource "google_cloud_run_service_iam_member" "dogs-pubsub" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.dogs.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.pubsub-pusher.email}"
}

resource "google_pubsub_subscription" "dogs" {
  project = local.project
  name    = "${local.prefix}-dogs-push"
  topic   = google_pubsub_topic.topic.name
  filter  = "hasPrefix(attributes.type, \"es.hotdoggi.events.dog_\")"
  push_config {
    push_endpoint = "${google_cloud_run_service.dogs.status[0].url}/events/"
    oidc_token {
      service_account_email = google_service_account.pubsub-pusher.email
    }
  }
}

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
    _ENVIRONMENT = "prod"
    _SERVICE     = "dogs"
    _REGION      = local.region
    _PREFIX      = local.prefix
  }
  filename = "services/dogs/cloudbuild.yaml"
}
