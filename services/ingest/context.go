package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	gin "github.com/gin-gonic/gin"
)

// Principal represents the the identity that originally authorized the context of an interaction
type Principal struct {
	ID         string `header:"id" firestore:"id" json:"id"`
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
	if Global["environment"].(string) != "prod" {
		devCaller := Principal{
			ID:         "1",
			Email:      "dev@localhost",
			Name:       "development",
			PictureURL: "unset",
		}
		c.Set("trace.id", "development mode")
		c.Set("principal", &devCaller)
		c.Set("principal.email", devCaller.Email)
		c.Set("principal.id", devCaller.ID)
		c.Set("principal.name", devCaller.Name)
		c.Set("principal.picture", devCaller.PictureURL)
		c.Next()
		return
	}

	// Trace context
	traceToken := c.GetHeader("X-Cloud-Trace-Context")
	// TODO check if correct format
	if traceToken != "" {
		c.Set("trace.id", traceToken)
	}

	encoded := c.Request.Header.Get("X-Endpoint-API-UserInfo")
	if encoded == "" {
		Respond(c, http.StatusUnauthorized, "missing gateway user info header")
		return
	}
	bytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		Respond(c, http.StatusUnauthorized, "failed to decode user info header")
		return
	}

	var apiCaller Principal
	err = json.Unmarshal(bytes, &apiCaller)
	if err != nil {
		Respond(c, http.StatusUnauthorized, "failed to deserialize user info header")
		return
	}

	// Context OK

	c.Set("principal", &apiCaller)
	c.Set("principal.email", apiCaller.Email)
	c.Set("principal.id", apiCaller.ID)
	c.Set("principal.name", apiCaller.Name)
	c.Set("principal.picture", apiCaller.PictureURL)
	c.Next()
}
