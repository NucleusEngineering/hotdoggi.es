# Copyright 2022 Daniel Stamer

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

resource "google_compute_managed_ssl_certificate" "frontend" {
  provider = google-beta
  project  = local.project
  name     = "${local.prefix}-certificate"

  managed {
    domains = [
      # "api.${local.domain}.",
      "${local.domain}."
    ]
  }
}

resource "google_compute_global_forwarding_rule" "frontend" {
  project               = local.project
  provider              = google-beta
  name                  = "${local.prefix}-forwarding-rule"
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL"
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
    hosts        = ["${local.domain}"]
    path_matcher = "app"
  }
  path_matcher {
    name            = "app"
    default_service = google_compute_backend_bucket.app.id
  }

}
