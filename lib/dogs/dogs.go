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

package dogs

import (
	"encoding/json"
	"fmt"
	"time"
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
	ID         string `header:"user_id" firestore:"user_id" json:"user_id"`
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

// Deserialize a ref into byte array
func (ref *DogRef) Deserialize(buffer []byte) error {
	err := json.Unmarshal(buffer, ref)
	if err != nil {
		return fmt.Errorf("failed deserialize payload: %v", err)
	}
	return nil
}

// Validate correctness of data
func (ref *DogRef) Validate(typeName string) error {
	switch typeName {
	case "es.hotdoggi.events.dog_added":
		// TODO implement type validation and return errors
		return nil
	case "es.hotdoggi.events.dog_removed":
		if ref.ID == "" {
			return fmt.Errorf("no reference id given for type: %s", typeName)
		}
		return nil
	case "es.hotdoggi.events.dog_updated":
		if ref.ID == "" {
			return fmt.Errorf("no reference id given for type: %s", typeName)
		}
		return nil
	case "es.hotdoggi.events.dog_moved":
		if ref.ID == "" {
			return fmt.Errorf("no reference id given for type: %s", typeName)
		}
		return nil
	default:
		return fmt.Errorf("unrecognized type name received: %s", typeName)
	}
}