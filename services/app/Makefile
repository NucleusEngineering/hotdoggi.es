SERVICE := app
PROJECT := hwsh-api
REGION := europe-west1

run:
	go run server.go

build:
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .

update:
	gcloud run services update ${PROJECT}-${SERVICE} --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE}

clean:
	rm -rf static
	rm -rf hwsh/node_modules

compile:
	cd hwsh && npm install
	cd hwsh && yarn build
	rm -rf static || exit 0
	mv hwsh/build static

.PHONY: run build update test static compile
