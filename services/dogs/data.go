package main

import (
	"context"
	"time"

	firestore "cloud.google.com/go/firestore"
	trace "go.opencensus.io/trace"
	iterator "google.golang.org/api/iterator"
)

const (
	collectionName = "es.hotdoggi.data.dogs"
)

// PubSubMessage is the data envelope used by pub/sub push subscriptions
type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// Event Data is the actual event 'data' payload
type EventData struct {
	Principal Principal `header:"principal" firestore:"principal" json:"principal"`
	Ref       DogRef    `header:"ref" firestore:"ref" json:"ref"`
}

// Principal represents the the identity that originally authorized the context of an interaction
type Principal struct {
	ID         string `header:"id" firestore:"id" json:"user_id"`
	Email      string `header:"email" firestore:"email" json:"email"`
	Name       string `header:"name" firestore:"name" json:"name"`
	PictureURL string `header:"picture" firestore:"picture" json:"picture"`
}

// DogRef is the actual content of the event data besides the calling principal
type DogRef struct {
	ID  string `header:"id" firestore:"id" json:"id"`
	Dog Dog    `header:"inline" firestore:"dog" json:"dog"`
}

// Dog data model
type Dog struct {
	Name       string   `header:"name" firestore:"name" json:"name"`
	Breed      string   `header:"breed" firestore:"breed" json:"breed"`
	Color      string   `header:"color" firestore:"color" json:"color"`
	Birthday   string   `header:"birthday" firestore:"birthday" json:"birthday"`
	PictureURL string   `header:"picture" firestore:"picture" json:"picture"`
	Location   Location `header:"inline" firestore:"location" json:"location"`
	Metadata   Metadata `header:"inline" firestore:"metadata" json:"metadata"`
}

// Location data model
type Location struct {
	Latitude  float32 `header:"latitude" firestore:"latitude" json:"latitude"`
	Longitude float32 `header:"longitude" firestore:"longitude" json:"longitude"`
}

// Metadata data model
type Metadata struct {
	Owner    string    `header:"owner" firestore:"owner" json:"owner"`
	Modified time.Time `firestore:"modified" json:"modified"`
}

// List all dogs
func List(ctx context.Context) ([]DogRef, error) {
	ctx, span := trace.StartSpan(ctx, "dogs.data.list")
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
func Get(ctx context.Context, key string) (DogRef, error) {
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
func Add(ctx context.Context, dog Dog) (DogRef, error) {
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
func Update(ctx context.Context, key string, dog Dog) (DogRef, error) {
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
func Delete(ctx context.Context, key string) error {
	ctx, span := trace.StartSpan(ctx, "dogs.data.delete")
	defer span.End()
	client := Global["client.firestore"].(*firestore.Client)
	_, err := client.Collection("dogs").Doc(key).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
