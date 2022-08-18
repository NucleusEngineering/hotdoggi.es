# Copyright 2022 Google

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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