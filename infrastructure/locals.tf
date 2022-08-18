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
