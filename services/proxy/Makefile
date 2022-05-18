SERVICE := proxy
PROJECT := $$(terraform output -json | jq -r '.project.value'))
REGION := $$(terraform output -json | jq -r '.region.value'))

init:
	go mod init

run:
	go run server.go

build:
	sh configure.sh
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .
	echo "undefined" > service.json

update:
	gcloud run services update ${PROJECT}-${SERVICE} --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE}

.PHONY: init run build update
