package main

import (
	"context"
	"time"

	firestore "cloud.google.com/go/firestore"
	gin "github.com/gin-gonic/gin"
	trace "go.opencensus.io/trace"
	iterator "google.golang.org/api/iterator"
)

const (
	collectionName = "es.hotdoggi.data.dogs"
)

// DogRef references a dog object
type DogRef struct {
	ID  string `header:"id" json:"id"`
	Dog Dog    `header:"inline" json:"dog"`
}

// Dog data model
type Dog struct {
	Name        string   `header:"name" firestore:"name" json:"name"`
	Description string   `header:"description" firestore:"description" json:"description"`
	Metadata    Metadata `header:"inline" firestore:"metadata" json:"metadata"`
}

// Metadata data model
type Metadata struct {
	Owner    string    `header:"owner" firestore:"owner" json:"owner"`
	Modified time.Time `firestore:"modified" json:"modified"`
}

// List all dogs
func List(ctx context.Context, c *gin.Context) ([]DogRef, error) {
	_, span := trace.StartSpan(ctx, "dogs.data.list")
	defer span.End()
	result := []DogRef{}
	client := Global["client.firestore"].(*firestore.Client)
	iter := client.Collection(collectionName).Documents(ctx)
	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var dog Dog
		snap.DataTo(&dog)
		result = append(result, DogRef{
			ID:  snap.Ref.ID,
			Dog: dog,
		})
	}

	return result, nil
}

// Get a specific dog
func Get(ctx context.Context, c *gin.Context, key string) (DogRef, error) {
	ctx, span := trace.StartSpan(ctx, "dogs.data.get")
	defer span.End()
	client := Global["client.firestore"].(*firestore.Client)
	snap, err := client.Collection(collectionName).Doc(key).Get(ctx)
	if err != nil {
		return DogRef{}, err
	}
	var dog Dog
	snap.DataTo(&dog)
	return DogRef{
		ID:  snap.Ref.ID,
		Dog: dog,
	}, nil
}

// Add a specific dog
func Add(ctx context.Context, c *gin.Context, dog Dog) (DogRef, error) {
	ctx, span := trace.StartSpan(ctx, "dogs.data.add")
	defer span.End()
	client := Global["client.firestore"].(*firestore.Client)

	dog.Metadata.Modified = time.Now()

	result, _, err := client.Collection(collectionName).Add(ctx, dog)
	if err != nil {
		return DogRef{}, err
	}
	return DogRef{
		ID:  result.ID,
		Dog: dog,
	}, nil
}

// Update a specific dog
func Update(ctx context.Context, c *gin.Context, key string, dog Dog) (DogRef, error) {
	ctx, span := trace.StartSpan(ctx, "dogs.data.update")
	defer span.End()
	client := Global["client.firestore"].(*firestore.Client)

	dog.Metadata.Modified = time.Now()

	_, err := client.Collection(collectionName).Doc(key).Set(ctx, dog)
	if err != nil {
		return DogRef{}, err
	}
	return DogRef{
		ID:  key,
		Dog: dog,
	}, nil
}

// Delete a specific dog
func Delete(ctx context.Context, c *gin.Context, key string) error {
	ctx, span := trace.StartSpan(ctx, "dogs.data.delete")
	defer span.End()
	client := Global["client.firestore"].(*firestore.Client)
	_, err := client.Collection("dogs").Doc(key).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
