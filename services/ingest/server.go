package main

import (
	"context"
	"log"
	"os"

	firestore "cloud.google.com/go/firestore"
	stackdriver "contrib.go.opencensus.io/exporter/stackdriver"
	gin "github.com/gin-gonic/gin"
	trace "go.opencensus.io/trace"
)

// Global map for shared resources
var Global map[string]interface{}

func main() {
	ctx := context.Background()
	configure(ctx)

	exporter := createTraceExporter()
	defer exporter.StopMetricsExporter()

	router := gin.Default()
	events := router.Group("/events")
	events.Use(UserContextFromAPI)
	{
		events.POST("/:type/:source", EventHandler)
	}
	router.Run()
}

func configure(ctx context.Context) {
	Global = make(map[string]interface{})
	Global["environment"] = os.Getenv("ENVIRONMENT")
	if Global["environment"].(string) == "" {
		Global["environment"] = "dev"
	}
	if Global["environment"].(string) == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	Global["project.id"] = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if Global["project.id"] == "" {
		log.Fatal("failed to read GOOGLE_CLOUD_PROJECT")
	}
	Global["client.firestore"] = createFirestoreClient(ctx)
}

func createTraceExporter() *stackdriver.Exporter {
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
