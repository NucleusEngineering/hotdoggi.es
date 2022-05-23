package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	firestore "cloud.google.com/go/firestore"
	stackdriver "contrib.go.opencensus.io/exporter/stackdriver"
	gin "github.com/gin-gonic/gin"
	trace "go.opencensus.io/trace"
)

func main() {
	ctx := context.Background()
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: project,
	})
	if err != nil {
		log.Fatalf("failed to initialize trace exporter: %v", err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(0.1)})

	client, err := firestore.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("failed to initialize firestore client: %v", err)
	}

	router := gin.Default()

	router.POST("/events/:type/:source", func(c *gin.Context) {
		ctx, span := trace.StartSpan(ctx, "ingest.handler.event")
		defer span.End()
		traceparent := fmt.Sprintf("00-%s-%s-0%d",
			span.SpanContext().TraceID.String(),
			span.SpanContext().SpanID.String(),
			span.SpanContext().TraceOptions,
		)
		log.Printf("tracep: %s\n", traceparent)
		code, err := validate(ctx, c)
		if err != nil {
			c.JSON(code, gin.H{
				"status": code,
				"error":  err.Error(),
			})
			return
		}
		code, err = commit(ctx, c, client, traceparent)
		if err != nil {
			c.JSON(code, gin.H{
				"status": code,
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "event inserted",
		})
	})

	router.Run()
}

func validate(ctx context.Context, c *gin.Context) (int, error) {
	_, span := trace.StartSpan(ctx, "ingest.validate")
	defer span.End()
	log.Printf("validating type\n")

	if c.Param("type") == "" {
		return http.StatusBadRequest, fmt.Errorf("no type name supplied")
	}

	if c.Param("source") == "" {
		return http.StatusBadRequest, fmt.Errorf("no source emitter supplied")
	}

	typeName := c.Param("type")

	switch typeName {
	case "es.hotdoggi.events.dog_added":
		// TODO implement type validation and return errors
		log.Printf("validation successful for type: %s\n", typeName)
	case "es.hotdoggi.events.dog_removed":
		// TODO implement type validation and return errors
		log.Printf("validation successful for type: %s\n", typeName)
	case "es.hotdoggi.events.dog_updated":
		// TODO implement type validation and return errors
		log.Printf("validation successful for type: %s\n", typeName)
	default:
		return http.StatusBadRequest, fmt.Errorf("unrecognized type name received: %s", typeName)
	}
	return 0, nil
}

func commit(ctx context.Context, c *gin.Context, client *firestore.Client, traceparent string) (int, error) {
	ctx, span := trace.StartSpan(ctx, "ingest.commit")
	defer span.End()
	typeName := c.Param("type")
	log.Printf("starting commit transaction for type: %s\n", typeName)

	buffer, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed read body payload: %v", err)
	}

	var obj map[string]interface{}
	err = json.Unmarshal(buffer, &obj)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed deserialize payload: %v", err)
	}

	_, _, err = client.Collection(typeName).Add(ctx, map[string]interface{}{
		"specversion":     "1.0",
		"subject":         "hotdoggi.es",
		"source":          c.Param("source"),
		"time":            firestore.ServerTimestamp,
		"traceparent":     traceparent,
		"datacontenttype": "application/json",
		"data":            obj,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to insert event into log: %v", err)
	}

	log.Printf("commit transaction successful for type: %s\n", typeName)
	return 0, nil
}
