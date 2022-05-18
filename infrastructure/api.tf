# Endpoints API service
resource "google_endpoints_service" "default" {
  service_name   = "api.${local.domain}"
  project        = local.project
  openapi_config = file("api.yaml")
}

output "config_id" {
  value = google_endpoints_service.default.config_id
}

output "service_name" {
  value = google_endpoints_service.default.service_name
}