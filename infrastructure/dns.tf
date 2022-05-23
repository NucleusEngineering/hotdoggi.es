resource "google_dns_managed_zone" "zone" {
  project     = local.project
  provider    = google-beta
  name        = local.prefix
  dns_name    = "${local.domain}."
  description = "Hotdoggies root zone"
}

resource "google_dns_record_set" "proxy-cname" {
  project     = local.project
  provider    = google-beta
  name         = "api.${google_dns_managed_zone.zone.dns_name}"
  managed_zone = google_dns_managed_zone.zone.name
  type         = "CNAME"
  ttl          = 300
  rrdatas      = ["ghs.googlehosted.com."]
}

resource "google_dns_record_set" "app-cname" {
  project     = local.project
  provider    = google-beta
  name         = "app.${google_dns_managed_zone.zone.dns_name}"
  managed_zone = google_dns_managed_zone.zone.name
  type         = "CNAME"
  ttl          = 300
  rrdatas      = ["ghs.googlehosted.com."]
}