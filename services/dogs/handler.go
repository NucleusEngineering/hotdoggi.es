package main

import (
	"context"
	"fmt"
	"net/http"

	gin "github.com/gin-gonic/gin"
	trace "go.opencensus.io/trace"
)

// ListHandler implements GET /
func ListHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	ctx, span := trace.StartSpan(ctx, "dogs.handler.list")
	defer span.End()

	result, err := List(ctx, c)
	if err != nil {
		Respond(c, http.StatusInternalServerError, fmt.Errorf("failed to retrieve objects: %v", err))
		return
	}
	Respond(c, http.StatusOK, result)
}

// GetHandler implements GET /{key}
func GetHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	ctx, span := trace.StartSpan(ctx, "dogs.handler.get")
	defer span.End()

	key := c.Param("key")
	result, err := Get(ctx, c, key)
	if err != nil {
		Respond(c, http.StatusInternalServerError, fmt.Errorf("failed to retrieve object: %v", err))
		return
	}
	Respond(c, http.StatusOK, result)
}

// EventHandler implements POST /
func EventHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	_, span := trace.StartSpan(ctx, "dogs.handler.event")
	defer span.End()

	// TODO implement
}

// func deserializeDog(c *gin.Context) (Dog, error) {
// 	body, err := ioutil.ReadAll(c.Request.Body)
// 	if err != nil {
// 		return Dog{}, err
// 	}

// 	var dog Dog
// 	err = json.Unmarshal(body, &dog)
// 	if err != nil {
// 		return Dog{}, err
// 	}

// 	return dog, nil
// }
