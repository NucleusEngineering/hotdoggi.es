SERVICE := proxy
PROJECT := $(shell cd ../../infrastructure && terraform output -json | jq -r '.project.value')
PREFIX := $(shell cd ../../infrastructure && terraform output -json | jq -r '.prefix.value')
REGION := $(shell cd ../../infrastructure && terraform output -json | jq -r '.region.value')
GATEWAY_SA := $(shell cd ../../infrastructure && terraform output -json | jq -r '.gateway_sa.value')

init:
	go mod init

run:
	go run server.go

build:
	gcloud services enable cloudbuild.googleapis.com
	sh configure.sh
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .
	echo "undefined" > service.json

update:
	gcloud run services update ${PREFIX}-${SERVICE} --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE}

.PHONY: init run build update
