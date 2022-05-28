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

	result, err := List(ctx)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": fmt.Sprintf("failed to retrieve objects: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetHandler implements GET /{key}
func GetHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	ctx, span := trace.StartSpan(ctx, "dogs.handler.get")
	defer span.End()

	key := c.Param("key")
	result, err := Get(ctx, key)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": fmt.Sprintf("failed to retrieve object: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// EventHandler implements POST /
func EventHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	ctx, span := trace.StartSpan(ctx, "dogs.handler.event")
	defer span.End()

	typeName := c.MustGet("event.type").(string)
	caller := c.MustGet("principal").(*Principal)
	ref := c.MustGet("event.data").(*DogRef)
	switch typeName {
	case "es.hotdoggi.events.dog_added":
		err := dogAdded(ctx, c, caller, ref)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to add dog: %v", err),
			})
		}
		c.Status(http.StatusCreated)
	case "es.hotdoggi.events.dog_removed":
		err := dogRemoved(ctx, c, caller, ref)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to remove dog: %v", err),
			})
		}
		c.Status(http.StatusAccepted)
	case "es.hotdoggi.events.dog_updated":
		err := dogUpdated(ctx, c, caller, ref)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to update dog: %v", err),
			})
		}
		c.Status(http.StatusAccepted)
	case "es.hotdoggi.events.dog_moved":
		err := dogUpdated(ctx, c, caller, ref)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to move dog: %v", err),
			})
		}
		c.Status(http.StatusAccepted)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("unknown event type received: %v", typeName),
		})
	}
}

func dogAdded(ctx context.Context, c *gin.Context, caller *Principal, ref *DogRef) error {
	ctx, span := trace.StartSpan(ctx, "dogs.handler.event.added")
	defer span.End()

	ref.Dog.Metadata.Owner = caller.ID
	_, err := Add(ctx, ref.Dog)
	return err
}

func dogRemoved(ctx context.Context, c *gin.Context, caller *Principal, ref *DogRef) error {
	ctx, span := trace.StartSpan(ctx, "dogs.handler.event.removed")
	defer span.End()
	existing, err := Get(ctx, ref.ID)
	if err != nil {
		return err
	}
	if existing.Dog.Metadata.Owner != caller.ID {
		return fmt.Errorf("refusing to remove dog. Dog not owned by caller")
	}

	err = Delete(ctx, ref.ID)
	return err
}

func dogUpdated(ctx context.Context, c *gin.Context, caller *Principal, ref *DogRef) error {
	ctx, span := trace.StartSpan(ctx, "dogs.handler.event.updated")
	defer span.End()
	existing, err := Get(ctx, ref.ID)
	if err != nil {
		return err
	}
	if existing.Dog.Metadata.Owner != caller.ID {
		return fmt.Errorf("refusing to update dog. Dog not owned by caller")
	}
	_, err = Update(ctx, ref.ID, ref.Dog)
	return err
}

func dogMoved(ctx context.Context, c *gin.Context, caller *Principal, id string, lat float32, long float32) error {
	ctx, span := trace.StartSpan(ctx, "dogs.handler.event.updated")
	defer span.End()
	existing, err := Get(ctx, id)
	if err != nil {
		return err
	}
	if existing.Dog.Metadata.Owner != caller.ID {
		return fmt.Errorf("refusing to move dog. Dog not owned by caller")
	}

	existing.Dog.Location.Latitude = lat
	existing.Dog.Location.Longitude = long

	_, err = Update(ctx, id, existing.Dog)
	return err
}
