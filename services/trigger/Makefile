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
