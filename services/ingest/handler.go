//  Copyright 2022 Google

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
	"io"
	"log"
	"net/http"

	firestore "cloud.google.com/go/firestore"
	gin "github.com/gin-gonic/gin"

	propagation "go.opentelemetry.io/otel/propagation"
	trace "go.opentelemetry.io/otel/trace"
)

// OptionsHandler for CORS preflights OPTIONS requests
func OptionsHandler(c *gin.Context) {}

// EventHandler implements POSTing events
func EventHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "ingest.handler:event")
	defer span.End()

	err := validate(ctx, c)
	if err != nil {
		log.Printf("error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to validate event: %v", err),
		})
		return
	}
	ref, err := commit(ctx, c)
	if err != nil {
		log.Printf("error: %v\n", err)
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
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	_, span := (*tracer).Start(ctx, "ingest.data:validate")
	defer span.End()

	sourceName := c.Param("source")
	if sourceName == "" {
		return fmt.Errorf("no source emitter supplied")
	}

	typeName := c.Param("type")
	if c.Param("type") == "" {
		return fmt.Errorf("no type name supplied")
	}

	buffer, err := io.ReadAll(c.Request.Body)
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
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "ingest.data:commit")
	defer span.End()

	client := Global["client.firestore"].(*firestore.Client)
	principal := c.MustGet("principal").(*Principal)
	typeName := c.MustGet("event.type").(string)
	sourceName := c.MustGet("event.source").(string)
	data := c.MustGet("event.data").(*DogRef)

	// Use explicit propagation of the trace context
	carrier := propagation.MapCarrier{}
	propagation.TraceContext{}.Inject(ctx, carrier)
	traceparent := carrier.Get("traceparent")

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
