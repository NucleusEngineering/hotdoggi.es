SERVICE := analytics
PROJECT := $$(terraform output -json | jq -r '.project.value')
REGION := $$(terraform output -json | jq -r '.region.value')
GATEWAY_SA := $$(terraform output -json | jq -r '.gateway_sa.value')

build:
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .

update:
	gcloud run services update ${PROJECT}-${SERVICE} --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE}

.PHONY: run build update
