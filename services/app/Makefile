SERVICE := app
PROJECT := $$(terraform output -json | jq -r '.project.value')
REGION := $$(terraform output -json | jq -r '.region.value')

run:
	go run server.go

build:
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .

update:
	gcloud run services update ${PROJECT}-${SERVICE} --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE}

clean:
	rm -rf static
	rm -rf spa/node_modules

compile:
	cd spa && npm install
	cd spa && yarn build
	rm -rf static || exit 0
	mv spa/build static

.PHONY: run build update test static compile
