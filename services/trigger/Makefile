SERVICE := trigger
PROJECT := $(shell cd ../../infrastructure && terraform output -json | jq -r '.project.value')
PREFIX := $(shell cd ../../infrastructure && terraform output -json | jq -r '.prefix.value')
BUCKET := $(shell cd ../../infrastructure && terraform output -json | jq -r '.function-bucket.value')
REGION := $(shell cd ../../infrastructure && terraform output -json | jq -r '.region.value')

build: clean
	zip -r function.zip *
	gsutil cp function.zip gs://${BUCKET}

clean:
	rm function.zip || exit 0

update:
	functions deploy ${PREFIX}-${SERVICE} --source=gsutil://${BUCKET}/function.zip --region=${REGION}

.PHONY: build clean update
