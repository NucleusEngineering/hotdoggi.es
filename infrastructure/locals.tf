locals {
  prefix         = "hotdoggies"
  region         = "europe-west1"
  project        = "hotdoggi-es"
  project_number = "640843850686"
  organization   = "780482598012"
  domain         = "hotdoggies.stamer.demo.altostrat.com"
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

output "prefix" {
  value = local.prefix
}
