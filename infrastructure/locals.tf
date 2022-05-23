locals {
  prefix         = "hotdoggies"
  region         = "europe-west4"
  project        = "hotdoggi-es"
  project_number = "640843850686"
  organization   = "780482598012"
  domain         = "hotdoggi.es"
  repo           = "hotdoggi.es"
  repo_owner     = "helloworlddan"
  branch         = "main"
  user           = "admin@stamer.altostrat.com"
}

output "project" {
  value = local.project
}
output "region" {
  value = local.region
}
