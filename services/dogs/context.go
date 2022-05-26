package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	gin "github.com/gin-gonic/gin"
)

// Principal represents the the identity that originally authorized the context of an interaction
type Principal struct {
	ID         string `header:"id" firestore:"id" json:"user_id"`
	Email      string `header:"email" firestore:"email" json:"email"`
	Name       string `header:"name" firestore:"name" json:"name"`
	PictureURL string `header:"picture" firestore:"picture" json:"picture"`
}

// UserContextFromAPI implements a middleware that resolves embedded user context info
// passed in from the firebase authentication at the service proxy layer.
func UserContextFromAPI(c *gin.Context) {
	ctx := c.Request.Context()
	c.Set("trace.context", ctx)

	// Skip verification in non-prod
	if Global["environment"].(string) == "dev" {
		devCaller := Principal{
			ID:         "1",
			Email:      "dev@localhost",
			Name:       "development",
			PictureURL: "unset",
		}
		c.Set("principal", &devCaller)
		c.Next()
		return
	}

	encoded := c.Request.Header.Get("X-Endpoint-API-UserInfo")
	if encoded == "" {
		c.JSON(http.StatusUnauthorized, "missing gateway user info header")
		c.Abort()
		return
	}
	bytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "failed to decode user info header")
		c.Abort()
		return
	}

	var apiCaller Principal
	err = json.Unmarshal(bytes, &apiCaller)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "failed to deserialize user info header")
		c.Abort()
		return
	}

	// Context OK
	c.Set("principal", &apiCaller)
	c.Next()
}

// UserContextFromEVENT implements a middleware that resolves embedded user context info
// passed in from the embedded event data.
func ContextFromEvent(c *gin.Context) {
	ctx := c.Request.Context()
	c.Set("trace.context", ctx)

	// Skip verification in non-prod
	if Global["environment"].(string) != "prod" {
		c.Next()
		return
	}

	// TODO build context injection
	c.Set("event.data", "uninjected ref object")

	// Trace context
	c.Set("trace.id", "uninjected")

	// Verification OK
	c.Set("principal", "uninjected pointer")
	c.Next()
}
