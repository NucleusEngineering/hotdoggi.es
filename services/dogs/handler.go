package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	gin "github.com/gin-gonic/gin"
	trace "go.opencensus.io/trace"
)

// ListHandler implements GET /
func ListHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	ctx, span := trace.StartSpan(ctx, "things.handler.list")
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
	ctx, span := trace.StartSpan(ctx, "things.handler.get")
	defer span.End()

	key := c.Param("key")
	result, err := Get(ctx, c, key)
	if err != nil {
		Respond(c, http.StatusInternalServerError, fmt.Errorf("failed to retrieve object: %v", err))
		return
	}
	Respond(c, http.StatusOK, result)
}

// PostHandler implements POST /
func PostHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	ctx, span := trace.StartSpan(ctx, "things.handler.post")
	defer span.End()

	thing, err := deserializeThing(c)
	if err != nil {
		Respond(c, http.StatusInternalServerError, fmt.Errorf("failed to deserialize payload: %v", err))
		return
	}

	result, err := Add(ctx, c, thing)
	if err != nil {
		Respond(c, http.StatusInternalServerError, fmt.Errorf("failed to store object: %v", err))
		return
	}
	c.JSON(http.StatusOK, result)
}

// PatchHandler implements PATCH /{key}
func PatchHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	ctx, span := trace.StartSpan(ctx, "things.handler.patch")
	defer span.End()

	key := c.Param("key")
	thing, err := deserializeThing(c)
	if err != nil {
		Respond(c, http.StatusInternalServerError, fmt.Errorf("failed to deserialize object: %v", err))
		return
	}

	result, err := Update(ctx, c, key, thing)
	if err != nil {
		Respond(c, http.StatusInternalServerError, fmt.Errorf("failed to update object: %v", err))
		return
	}
	Respond(c, http.StatusOK, result)
}

// DeleteHandler implements DELETE /{key}
func DeleteHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	ctx, span := trace.StartSpan(ctx, "things.handler.delete")
	defer span.End()

	key := c.Param("key")
	err := Delete(ctx, c, key)
	if err != nil {
		Respond(c, http.StatusInternalServerError, fmt.Errorf("failed to delete object: %v", err))
		return
	}
	Respond(c, http.StatusOK, nil)
}

func deserializeThing(c *gin.Context) (Thing, error) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return Thing{}, err
	}

	var thing Thing
	err = json.Unmarshal(body, &thing)
	if err != nil {
		return Thing{}, err
	}

	return thing, nil
}
