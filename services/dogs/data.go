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
	"fmt"
	"log"
	"time"

	firestore "cloud.google.com/go/firestore"
	trace "go.opentelemetry.io/otel/trace"
	iterator "google.golang.org/api/iterator"

	dogs "github.com/helloworlddan/hotdoggi.es/lib/dogs"
)

const (
	collectionName = "es.hotdoggi.data.dogs"
)

// List all dogs
func List(ctx context.Context, userID string) ([]dogs.DogRef, error) {
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.data:list")
	defer span.End()

	result := []dogs.DogRef{}
	client := Global["client.firestore"].(*firestore.Client)
	log.Printf("scanning dogs\n")
	iter := client.Collection(collectionName).Where("metadata.owner", "==", userID).Documents(ctx)
	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var dog dogs.Dog
		snap.DataTo(&dog)
		result = append(result, dogs.DogRef{
			ID:  snap.Ref.ID,
			Dog: dog,
		})
	}

	return result, nil
}

// ListStream will stream updates for all dogs belonging to a user back on a channel until a quit message is sent.
func ListStream(ctx context.Context, userID string, dogChan chan<- dogs.DogRef) {
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.data:list#streaming")
	defer span.End()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := Global["client.firestore"].(*firestore.Client)

	log.Printf("streaming dogs for user(%s)\n", userID)

	iter := client.Collection(collectionName).Where("metadata.owner", "==", userID).Snapshots(ctx)
	for {
		snap, err := iter.Next() // this is the blocking listener
		if err != nil {
			close(dogChan)
			return
		}
		if snap != nil {
			for _, change := range snap.Changes {
				switch change.Kind {
				case firestore.DocumentAdded:
					// TODO what to do with different types of events?
					// -> Nothing for now
				case firestore.DocumentRemoved:
					// TODO what to do with different types of events?
					// -> Nothing for now
				case firestore.DocumentModified:
					var dog dogs.Dog
					change.Doc.DataTo(&dog)

					dogChan <- dogs.DogRef{
						ID:  change.Doc.Ref.ID,
						Dog: dog,
					}
				}
			}
		}
	}
}

// Get a specific dog
func Get(ctx context.Context, userID string, key string) (dogs.DogRef, error) {
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.data:get")
	defer span.End()

	client := Global["client.firestore"].(*firestore.Client)

	log.Printf("reading dog(%s)\n", key)
	snap, err := client.Collection(collectionName).Doc(key).Get(ctx)
	if err != nil {
		return dogs.DogRef{}, err
	}
	var dog dogs.Dog
	snap.DataTo(&dog)

	if dog.Metadata.Owner != userID {
		return dogs.DogRef{}, fmt.Errorf("dog not owned by user %v", userID)
	}

	return dogs.DogRef{
		ID:  snap.Ref.ID,
		Dog: dog,
	}, nil
}

// GetStream will stream updates for a specific dog back on a channel until a quit message is sent.
func GetStream(ctx context.Context, userID string, key string, dogChan chan<- dogs.DogRef) {
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.data:get#streaming")
	defer span.End()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := Global["client.firestore"].(*firestore.Client)

	log.Printf("streaming updates for dog(%s)\n", userID)

	iter := client.Collection(collectionName).Doc(key).Snapshots(ctx)
	for {
		snap, err := iter.Next() // this is the blocking listener
		if err != nil {
			close(dogChan)
			return
		}
		if snap != nil {
			var dog dogs.Dog
			snap.DataTo(&dog)
			dogRef := dogs.DogRef{
				ID:  snap.Ref.ID,
				Dog: dog,
			}

			if dogRef.Dog.Metadata.Owner != userID {
				close(dogChan)
				return
			}
			dogChan <- dogRef
		}
	}
}

// Add a specific dog
func Add(ctx context.Context, dog dogs.Dog) (dogs.DogRef, error) {
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.data:add")
	defer span.End()

	client := Global["client.firestore"].(*firestore.Client)

	dog.Metadata.Modified = time.Now()

	log.Printf("adding new dog ... ")
	ref, _, err := client.Collection(collectionName).Add(ctx, dog)
	if err != nil {
		return dogs.DogRef{}, err
	}

	return dogs.DogRef{
		ID:  ref.ID,
		Dog: dog,
	}, nil
}

// Update a specific dog
func Update(ctx context.Context, key string, dog dogs.Dog) (dogs.DogRef, error) {
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.data:update")
	defer span.End()

	client := Global["client.firestore"].(*firestore.Client)

	dog.Metadata.Modified = time.Now()

	log.Printf("updating dog(%s)\n", key)
	_, err := client.Collection(collectionName).Doc(key).Set(ctx, dog)
	if err != nil {
		return dogs.DogRef{}, err
	}

	return dogs.DogRef{
		ID:  key,
		Dog: dog,
	}, nil
}

// Delete a specific dog
func Delete(ctx context.Context, key string) error {
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.data:delete")
	defer span.End()

	client := Global["client.firestore"].(*firestore.Client)

	log.Printf("deleting dog(%s)\n", key)
	_, err := client.Collection(collectionName).Doc(key).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
