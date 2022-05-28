package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	firestore "cloud.google.com/go/firestore"
	gin "github.com/gin-gonic/gin"
	trace "go.opencensus.io/trace"
)


// EventHandler implements POSTing events
func EventHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	ctx, span := trace.StartSpan(ctx, "ingest.handler.event")
	defer span.End()

	err := validate(ctx, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to validate event: %v", err),
		})
		return
	}
	ref, err := commit(ctx, c)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": fmt.Sprintf("failed to commit to event log: %v", err),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":  "event inserted",
		"event_id": ref.ID,
	})
}

func validate(ctx context.Context, c *gin.Context) error {
	_, span := trace.StartSpan(ctx, "ingest.validate")
	defer span.End()

	sourceName := c.Param("source")
	if sourceName == "" {
		return fmt.Errorf("no source emitter supplied")
	}

	typeName := c.Param("type")
	if c.Param("type") == "" {
		return fmt.Errorf("no type name supplied")
	}

	buffer, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Errorf("failed read body payload: %v", err)
	}

	// Assume we are talking about dogs
	var data DogRef
	err = (&data).deserialize(buffer)
	if err != nil {
		return fmt.Errorf("failed deserialize payload: %v", err)
	}
	err = (&data).validate(typeName)
	if err != nil {
		return fmt.Errorf("failed validate payload: %v", err)
	}

	// Validation OK
	c.Set("event.data", &data)
	c.Set("event.type", typeName)
	c.Set("event.source", sourceName)

	return nil
}

func commit(ctx context.Context, c *gin.Context) (*firestore.DocumentRef, error) {
	ctx, span := trace.StartSpan(ctx, "ingest.commit")
	defer span.End()

	client := Global["client.firestore"].(*firestore.Client)
	principal := c.MustGet("principal").(*Principal)
	typeName := c.MustGet("event.type").(string)
	sourceName := c.MustGet("event.source").(string)
	data := c.MustGet("event.data").(*DogRef)

	traceparent := fmt.Sprintf("00-%s-%s-0%d",
		span.SpanContext().TraceID.String(),
		span.SpanContext().SpanID.String(),
		span.SpanContext().TraceOptions,
	)

	payload := EventData{
		Principal: *principal,
		Ref:       *data,
	}

	ref, _, err := client.Collection(typeName).Add(ctx, map[string]interface{}{
		"specversion":     "1.0",
		"subject":         "hotdoggi.es",
		"source":          sourceName,
		"time":            firestore.ServerTimestamp,
		"traceparent":     traceparent,
		"datacontenttype": "application/json",
		"data":            payload,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert event: %v", err)
	}

	return ref, nil
}
