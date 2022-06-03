resource "google_storage_bucket" "app" {
  project                     = local.project
  provider                    = google-beta
  name                        = "${local.prefix}-app-bucket"
  uniform_bucket_level_access = true
  location                    = "EU"

  force_destroy = true
  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
}

resource "google_compute_backend_bucket" "app" {
  project     = local.project
  provider    = google-beta
  name        = "${local.prefix}-app-backend"
  description = "Origin for static SPA serving"
  bucket_name = google_storage_bucket.app.name
  enable_cdn  = true
}

resource "google_storage_bucket_iam_member" "app" {
  bucket = google_storage_bucket.app.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

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
  description = "Build pipeline for ${local.prefix}-app"
  substitutions = {
    _REGION  = local.region
    _PREFIX  = local.prefix
    _SERVICE = "app"
    _BUCKET  = google_storage_bucket.app.name
  }

  filename = "app/cloudbuild.yaml"
}
