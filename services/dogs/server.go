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

package main

import (
	"context"
	"log"
	"net/http"
	"os"

	firestore "cloud.google.com/go/firestore"
	stackdriver "contrib.go.opencensus.io/exporter/stackdriver"
	firebase "firebase.google.com/go/v4"
	auth "firebase.google.com/go/v4/auth"
	gin "github.com/gin-gonic/gin"
	ochttp "go.opencensus.io/plugin/ochttp"
	propagation "go.opencensus.io/plugin/ochttp/propagation/b3"
	trace "go.opencensus.io/trace"
)

// Global map for shared resources
var Global map[string]interface{}

func main() {
	ctx := context.Background()
	configure(ctx)

	exporter := Global["client.trace"].(*stackdriver.Exporter)
	defer exporter.StopMetricsExporter()

	router := gin.Default()
	events := router.Group("/events")
	events.Use(ContextFromEvent)
	{
		events.POST("/", EventHandler)
	}

	api := router.Group("/dogs")
	api.Use(UserContextFromAPI)
	{
		api.GET("/:key", GetHandler)
		api.GET("/", ListHandler)

	}
	router.Run()
}

func configure(ctx context.Context) {
	Global = make(map[string]interface{})
	Global["environment"] = os.Getenv("ENVIRONMENT")

	// Default to prod (safer)
	if Global["environment"].(string) == "" {
		Global["environment"] = "prod"
	}

	if Global["environment"] == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	Global["project.id"] = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if Global["project.id"] == "" {
		log.Fatal("failed to read GOOGLE_CLOUD_PROJECT")
	}
	Global["client.http"] = createHTTPClient(ctx)
	Global["client.firebase"] = createFirebaseClient(ctx)
	Global["client.firestore"] = createFirestoreClient(ctx)
	Global["client.trace"] = createTraceExporter(ctx)
}

func createTraceExporter(ctx context.Context) *stackdriver.Exporter {
	projectID := Global["project.id"].(string)
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(0.1)})
	exporter.StartMetricsExporter()
	return exporter
}

func createFirestoreClient(ctx context.Context) *firestore.Client {
	projectID := Global["project.id"].(string)
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	return client
}

func createFirebaseClient(ctx context.Context) *auth.Client {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("failed to create firebase app: %v", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("failed to create firebase client: %v", err)
	}
	return client
}

func createHTTPClient(ctx context.Context) *http.Client {
	client := &http.Client{
		Transport: &ochttp.Transport{
			Propagation: &propagation.HTTPFormat{},
		},
	}
	return client
}
