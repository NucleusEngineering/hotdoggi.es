terraform {
  backend "gcs" {
    bucket = "hotdoggi-es-terraform-state"
    prefix = "terraform/state"
  }
}

provider "google-beta" {
  region = local.region
}

provider "google" {
  region = local.region
}
