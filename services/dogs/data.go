package main

import (
	"context"
	"time"

	firestore "cloud.google.com/go/firestore"
	gin "github.com/gin-gonic/gin"
	trace "go.opencensus.io/trace"
	iterator "google.golang.org/api/iterator"
)

// ThingRef references a thing object
type ThingRef struct {
	ID    string `header:"id" json:"id"`
	Thing Thing  `header:"inline" json:"thing"`
}

// Thing data model
type Thing struct {
	Name        string   `header:"name" firestore:"name" json:"name"`
	Description string   `header:"description" firestore:"description" json:"description"`
	Metadata    Metadata `header:"inline" firestore:"metadata" json:"metadata"`
}

// Metadata data model
type Metadata struct {
	Owner    string    `header:"owner" firestore:"owner"  json:"owner"`
	Modified time.Time `firestore:"modified" json:"modified"`
}

// List all things
func List(ctx context.Context, c *gin.Context) ([]ThingRef, error) {
	_, span := trace.StartSpan(ctx, "things.data.list")
	defer span.End()
	result := []ThingRef{}
	client := Global["client.firestore"].(*firestore.Client)
	iter := client.Collection("things").Documents(ctx)
	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var thing Thing
		snap.DataTo(&thing)
		result = append(result, ThingRef{
			ID:    snap.Ref.ID,
			Thing: thing,
		})
	}

	return result, nil
}

// Get a specific thing
func Get(ctx context.Context, c *gin.Context, key string) (ThingRef, error) {
	ctx, span := trace.StartSpan(ctx, "things.data.get")
	defer span.End()
	client := Global["client.firestore"].(*firestore.Client)
	snap, err := client.Collection("things").Doc(key).Get(ctx)
	if err != nil {
		return ThingRef{}, err
	}
	var thing Thing
	snap.DataTo(&thing)
	return ThingRef{
		ID:    snap.Ref.ID,
		Thing: thing,
	}, nil
}

// Add a specific thing
func Add(ctx context.Context, c *gin.Context, thing Thing) (ThingRef, error) {
	ctx, span := trace.StartSpan(ctx, "things.data.add")
	defer span.End()
	client := Global["client.firestore"].(*firestore.Client)

	thing.Metadata.Owner = c.MustGet("caller.id").(string)
	thing.Metadata.Modified = time.Now()

	result, _, err := client.Collection("things").Add(ctx, thing)
	if err != nil {
		return ThingRef{}, err
	}
	return ThingRef{
		ID:    result.ID,
		Thing: thing,
	}, nil
}

// Update a specific thing
func Update(ctx context.Context, c *gin.Context, key string, thing Thing) (ThingRef, error) {
	ctx, span := trace.StartSpan(ctx, "things.data.update")
	defer span.End()
	client := Global["client.firestore"].(*firestore.Client)

	thing.Metadata.Owner = c.MustGet("caller.id").(string)
	thing.Metadata.Modified = time.Now()

	_, err := client.Collection("things").Doc(key).Set(ctx, thing)
	if err != nil {
		return ThingRef{}, err
	}
	return ThingRef{
		ID:    key,
		Thing: thing,
	}, nil
}

// Delete a specific thing
func Delete(ctx context.Context, c *gin.Context, key string) error {
	ctx, span := trace.StartSpan(ctx, "things.data.delete")
	defer span.End()
	client := Global["client.firestore"].(*firestore.Client)
	_, err := client.Collection("things").Doc(key).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
