package main

import (
	"encoding/json"
	"fmt"
)

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
}

// Location data model
type Location struct {
	Latitude  float32 `header:"latitude" firestore:"latitude" json:"latitude"`
	Longitude float32 `header:"longitude" firestore:"longitude" json:"longitude"`
}

func (ref *DogRef) deserialize(buffer []byte) error {
	err := json.Unmarshal(buffer, ref)
	if err != nil {
		return fmt.Errorf("failed deserialize payload: %v", err)
	}
	return nil
}

func (ref *DogRef) validate(typeName string) error {
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
		if ref.Dog.Location.Latitude == 0.0 || ref.Dog.Location.Longitude == 0.0 {
			return fmt.Errorf("no positional coordinates give for : %s", typeName)
		}
		return nil
	default:
		return fmt.Errorf("unrecognized type name received: %s", typeName)
	}
}
