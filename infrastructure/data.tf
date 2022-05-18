# Firestore database for this project
resource "google_app_engine_application" "database" {
  project       = local.project
  provider      = google-beta
  location_id   = replace(replace(replace(replace(local.region, "1", ""), "2", ""), "3", ""), "4", "") # TODO strip all trailing numbers
  database_type = "CLOUD_FIRESTORE"
}
