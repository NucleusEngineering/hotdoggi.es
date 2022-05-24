SERVICE := analytics
PROJECT := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.project.value')
PREFIX := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.prefix.value')
REGION := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.region.value')
GATEWAY_SA := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.gateway_sa.value')

build:
	gcloud services enable cloudbuild.googleapis.com
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .

update:
	gcloud run services update ${PREFIX}-${SERVICE} --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE}

.PHONY: build update
