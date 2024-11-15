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
	"log"
	"net/http"

	gin "github.com/gin-gonic/gin"
	dogs "github.com/helloworlddan/hotdoggi.es/lib/dogs"
	trace "go.opentelemetry.io/otel/trace"
)

// UserContextFromAPI implements a middleware that resolves embedded user context info
// passed in from firebase authentication at the service proxy layer.
func UserContextFromAPI(c *gin.Context) {
	// Support CORS preflight OPTIONS requests
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	// Explicitly create new context, start of trace
	ctx := context.Background()
	// Create root span
	ctx, span := (*tracer).Start(ctx, "event", trace.WithNewRoot())
	span.End()

	ctx, span = (*tracer).Start(ctx, "ingest.context:api")
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

	c.Set("principal", &caller)

	// Context OK
	c.Next()
}
