resource "google_compute_managed_ssl_certificate" "frontend" {
  provider = google-beta
  project  = local.project
  name     = "${local.prefix}-certificate"

  managed {
    domains = [
      "${local.domain}.",
      "api.${local.domain}."
    ]
  }
}

resource "google_compute_global_forwarding_rule" "frontend" {
  project               = local.project
  provider              = google-beta
  name                  = "${local.prefix}-forwarding-rule"
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  port_range            = "443"
  target                = google_compute_target_https_proxy.frontend.id
}

resource "google_compute_target_https_proxy" "frontend" {
  project          = local.project
  provider         = google-beta
  name             = "${local.prefix}-https-proxy"
  url_map          = google_compute_url_map.frontend.id
  ssl_certificates = [google_compute_managed_ssl_certificate.frontend.id]
}

resource "google_compute_url_map" "frontend" {
  project         = local.project
  provider        = google-beta
  name            = "${local.prefix}-urlmap"
  default_service = google_compute_backend_bucket.app.id

  host_rule {
    hosts        = ["api.${local.domain}"]
    path_matcher = "api"
  }
  path_matcher {
    name            = "api"
    default_service = google_compute_backend_service.proxy.id
  }
  host_rule {
    hosts        = ["${local.domain}"]
    path_matcher = "app"
  }
  path_matcher {
    name            = "app"
    default_service = google_compute_backend_bucket.app.id
  }

}
