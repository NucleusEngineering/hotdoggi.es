//  Copyright 2022 Daniel Stamer

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at

//      http://www.apache.org/licenses/LICENSE-2.0

//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

module main

go 1.16

require (
	cloud.google.com/go/firestore v1.5.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.6
	github.com/gin-gonic/gin v1.7.1
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	go.opencensus.io v0.23.0
)
