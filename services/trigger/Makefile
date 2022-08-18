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

SERVICE := trigger
PROJECT := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.project.value')
PREFIX := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.prefix.value')
BUCKET := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.function_bucket.value')
REGION := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.region.value')

build: clean
	zip -r function.zip *
	gsutil cp function.zip gs://${BUCKET}

clean:
	rm -rf vendor/ || exit 0
	rm function.zip || exit 0

init:
	bundle install

run: init
	bundle exec functions-framework-ruby --target function

test: run clean

update:
	functions deploy ${PREFIX}-${SERVICE} --source=gsutil://${BUCKET}/function.zip --region=${REGION}

.PHONY: build clean update test init run
