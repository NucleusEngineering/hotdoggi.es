# Copyright 2022 Google

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

all: analytics archiver ingest services app

deploy:
	make -C infrastructure init apply

destroy:
	make -C infrastructure destroy

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

client:
	make -C clients/dogs install

load:
	@clear
	python3 loader/loader.py

services: dogs

.PHONY: static deploy destroy proxy services dogs analytics archiver ingest trigger client
