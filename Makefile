all: analytics archiver ingest services app

deploy:
	make -C infrastructure init apply

destroy:
	make -C infrastructure destroy

app:
	make -C services/app compile build clean

proxy:
	make -C services/proxy build

dogs:
	make -C services/dogs build

analytics:
	make -C services/analytics build

archiver:
	make -C services/archiver build

ingest:
	make -C services/ingest build

trigger:
	make -C services/trigger build

services: dogs

.PHONY: static deploy destroy app proxy services dogs analytics archiver ingest trigger
