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
	"fmt"
	"log"
	"os"

	firestore "cloud.google.com/go/firestore"
	gin "github.com/gin-gonic/gin"

	exporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	otel "go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	trace "go.opentelemetry.io/otel/trace"
)

const (
	prefixIdentifier string = "es.hotdoggi"
	serviceName      string = "ingest"
)

// Global map for shared resources
var Global map[string]interface{}

func main() {
	ctx := context.Background()
	configure(ctx)

	provider := Global["client.trace.provider"].(*sdktrace.TracerProvider)
	defer provider.ForceFlush(ctx)

	router := gin.Default()
	events := router.Group("/v1/events")
	events.Use(UserContextFromAPI)
	{
		events.POST("/:type/:source", EventHandler)
	}
	log.Println("Starting server.")
	log.Fatalf("error: %v", router.Run())
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
	Global["client.firestore"] = createFirestoreClient(ctx)
	Global["client.trace.exporter"] = createTraceExporter(ctx)
	Global["client.trace.provider"] = createTraceProvider(ctx)
	Global["client.trace.tracer"] = createTracer(ctx)
	log.Println("Configuration complete.")
}

func createTraceExporter(ctx context.Context) *exporter.Exporter {
	projectID := Global["project.id"].(string)
	exporter, err := exporter.New(exporter.WithProjectID(projectID))
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}

	return exporter
}

func createTraceProvider(ctx context.Context) *sdktrace.TracerProvider {
	exporter := Global["client.trace.exporter"].(*exporter.Exporter)

	// Probabilistic trace exporter in PROD
	provider := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)))
	// AlwaysOn trace exporter in DEV
	if Global["environment"].(string) == "dev" {
		provider = sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	}
	otel.SetTracerProvider(provider)
	return provider
}

func createTracer(ctx context.Context) *trace.Tracer {
	tracer := otel.GetTracerProvider().Tracer(fmt.Sprintf("%s.service.%s/", prefixIdentifier, serviceName))
	return &tracer
}

func createFirestoreClient(ctx context.Context) *firestore.Client {
	projectID := Global["project.id"].(string)
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	return client
}
