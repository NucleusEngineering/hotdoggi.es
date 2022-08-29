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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	gin "github.com/gin-gonic/gin"
	propagation "go.opentelemetry.io/otel/propagation"
	trace "go.opentelemetry.io/otel/trace"

	dogs "github.com/helloworlddan/hotdoggi.es/lib/dogs"
)

// UserContextFromAPI implements a middleware that resolves embedded user context info
// passed in from firebase authentication at the service proxy layer.
func UserContextFromAPI(c *gin.Context) {
	// Support CORS preflight OPTIONS requests
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	// Explicitly create new context, start of trace
	ctx := context.Background()
	// Create root span
	ctx, span := (*tracer).Start(ctx, "api", trace.WithNewRoot())
	span.End()

	ctx, span = (*tracer).Start(ctx, "dogs.context:api")
	defer span.End()
	c.Set("trace.context", ctx)

	// Skip verification in non-prod
	if Global["environment"].(string) == "dev" {
		caller := dogs.Principal{
			ID:         "1",
			Email:      "dev@localhost",
			Name:       "development",
			PictureURL: "unset",
		}
		c.Set("principal", &caller)
		c.Next()
		return
	}

	encoded := c.Request.Header.Get("X-Endpoint-API-UserInfo")
	if encoded == "" {
		log.Printf("error: %v\n", fmt.Errorf("missing gateway user info header"))
		c.JSON(http.StatusUnauthorized, "missing gateway user info header")
		c.Abort()
		return
	}
	bytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		log.Printf("error: %v\n", err)
		c.JSON(http.StatusUnauthorized, "failed to decode user info header")
		c.Abort()
		return
	}

	var caller dogs.Principal
	err = json.Unmarshal(bytes, &caller)
	if err != nil {
		log.Printf("error: %v\n", err)
		c.JSON(http.StatusUnauthorized, "failed to deserialize user info header")
		c.Abort()
		return
	}

	// Context OK
	c.Set("principal", &caller)
	c.Next()
}

// UserContextFromEvent implements a middleware that resolves embedded user context info
// passed in from the event data.
func ContextFromEvent(c *gin.Context) {
	buffer, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error: %v\n", err)
		c.JSON(http.StatusBadRequest, fmt.Errorf("failed to read POST payload: %v", err))
		c.Abort()
		return
	}

	var psMessage dogs.PubSubMessage
	err = json.Unmarshal(buffer, &psMessage)
	if err != nil {
		log.Printf("error: %v\n", err)
		c.JSON(http.StatusBadRequest, fmt.Errorf("failed deserialize Pub/Sub envelope: %v", err))
		c.Abort()
		return
	}

	event := cloudevents.NewEvent()
	err = json.Unmarshal(psMessage.Message.Data, &event)
	if err != nil {
		log.Printf("error: %v\n", err)
		c.JSON(http.StatusNotAcceptable, fmt.Errorf("failed to deserialize CloudEvent payload: %v", err))
		c.Abort()
		return
	}

	c.Set("event.type", event.Context.GetType())
	c.Set("event.source", event.Context.GetSource())

	// Resume trace
	// Explicitly override context from original event trace
	traceparent, err := event.Context.GetExtension("traceparent")
	if err != nil {
		log.Printf("error: %v\n", err)
		c.JSON(http.StatusNotAcceptable, fmt.Errorf("failed to read traceparent from event context: %v", err))
		c.Abort()
		return
	}
	ctx := propagation.TraceContext{}.Extract(context.Background(), propagation.MapCarrier{
		"traceparent": traceparent.(string),
	})
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.context:event")
	defer span.End()
	c.Set("trace.context", ctx)

	var data dogs.EventData
	err = json.Unmarshal(event.Data(), &data)
	if err != nil {
		log.Printf("error: %v\n", err)
		c.JSON(http.StatusNotAcceptable, fmt.Errorf("failed to deserialize resources: %v", err))
		c.Abort()
		return
	}

	c.Set("event.data", &data.Ref)
	c.Set("principal", &data.Principal)

	log.Printf("received event %s of type %s from source %s in context of user %s", event.ID(), event.Type(), event.Source(), data.Principal.Email)

	// Context OK
	c.Next()
}
