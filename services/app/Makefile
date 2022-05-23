SERVICE := app
PROJECT := $(shell cd ../../infrastructure && terraform output -json | jq -r '.project.value')
PREFIX := $(shell cd ../../infrastructure && terraform output -json | jq -r '.prefix.value')
REGION := $(shell cd ../../infrastructure && terraform output -json | jq -r '.region.value')
GATEWAY_SA := $(shell cd ../../infrastructure && terraform output -json | jq -r '.gateway_sa.value')

run:
	go run server.go

build:
	gcloud services enable cloudbuild.googleapis.com
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .

update:
	gcloud run services update ${PREFIX}-${SERVICE} --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE}

clean:
	rm -rf static
	rm -rf spa/node_modules

compile:
	cd spa && npm install
	cd spa && yarn build
	rm -rf static || exit 0
	mv spa/build static

.PHONY: run build update clean compile
