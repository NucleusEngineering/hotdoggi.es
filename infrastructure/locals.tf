locals {
  prefix         = "hotdoggies"
  region         = "europe-west4"
  project        = "hotdoggi-es"
  project_number = "" # TODO replace
  organization   = "" # TODO replace
  domain         = "" # TODO replace
  repo           = "hotdoggi.es"
  repo_owner     = "helloworlddan"
  branch         = "main"
}

output "project" {
  value = local.project
}
output "region" {
  value = local.region
}
