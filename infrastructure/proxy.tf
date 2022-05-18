# Disable authentication check to invoke this service
resource "google_cloud_run_service_iam_binding" "public" {
  location = local.region
  project  = local.project
  service  = google_cloud_run_service.proxy.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
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

# Service definition
resource "google_cloud_run_service" "proxy" {
  provider = google-beta
  project  = local.project
  name     = "${local.prefix}-proxy"
  location = local.region
  template {
    spec {
      service_account_name = google_service_account.proxy.email
      containers {
        image = "gcr.io/${local.project}/proxy"
      }
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
resource "google_service_account_iam_binding" "proxy-sa-user" {
  service_account_id = google_service_account.proxy.name
  role               = "roles/iam.serviceAccountUser"
  members = [
    "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
  ]
}

output "gateway_sa" {
  value = google_service_account.proxy.email
} 

