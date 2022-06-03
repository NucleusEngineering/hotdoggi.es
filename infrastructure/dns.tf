resource "google_dns_managed_zone" "zone" {
  project     = local.project
  provider    = google-beta
  name        = local.prefix
  dns_name    = "${local.domain}."
  description = "Hotdoggies root zone"
}

resource "google_dns_record_set" "lb-a-app" {
  project      = local.project
  provider     = google-beta
  name         = "${local.domain}."
  managed_zone = google_dns_managed_zone.zone.name
  type         = "A"
  ttl          = 300
  rrdatas      = [google_compute_global_forwarding_rule.frontend.ip_address]
}

resource "google_dns_record_set" "lb-a-api" {
  project      = local.project
  provider     = google-beta
  name         = "api.${local.domain}."
  managed_zone = google_dns_managed_zone.zone.name
  type         = "CNAME"
  ttl          = 300
  rrdatas      = ["ghs.googlehosted.com."]
}