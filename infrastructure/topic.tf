resource "google_pubsub_topic" "topic" {
  project = local.project
  name    = "${local.prefix}-stream"
}

resource "google_pubsub_topic" "dead-letter" {
  project = local.project
  name    = "${local.prefix}-dead-letters"
}

resource "google_pubsub_subscription" "subscription" {
  project = local.project
  name    = "${local.prefix}-dead-letter-pull"
  topic   = google_pubsub_topic.dead-letter.name
}

resource "google_project_iam_member" "pubsub-sa-tokencreator" {
  project = local.project
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:service-${local.project_number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "pubsub-sa-publisher" {
  project = local.project
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:service-${local.project_number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_project_iam_member" "pubsub-sa-subscriber" {
  project = local.project
  role    = "roles/pubsub.subscriber"
  member  = "serviceAccount:service-${local.project_number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_service_account" "pubsub-pusher" {
  project      = local.project
  account_id   = "${local.prefix}-pubsub-pusher"
  display_name = "${local.prefix}-pubsub-pusher"
}