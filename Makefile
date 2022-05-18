all:
	app dogs proxy

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

services: dogs

.PHONY: static deploy destroy app proxy services dogs
