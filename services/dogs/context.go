package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	gin "github.com/gin-gonic/gin"
)

// UserContextFromAPI implements a middleware that resolves embedded user context info
// passed in from firebase authentication at the service proxy layer.
func UserContextFromAPI(c *gin.Context) {
	ctx := c.Request.Context()
	c.Set("trace.context", ctx)

	// Skip verification in non-prod
	if Global["environment"].(string) == "dev" {
		caller := Principal{
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

	var caller Principal
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
	ctx := c.Request.Context()
	c.Set("trace.context", ctx)

	buffer, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error: %v\n", err)
		c.JSON(http.StatusBadRequest, fmt.Errorf("failed to read POST payload: %v", err))
		c.Abort()
		return
	}

	var psMessage PubSubMessage
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

	var data EventData
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
