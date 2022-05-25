package main

import (
	"encoding/json"
	"fmt"
)

// DogRef is the actual content of the event data besides the calling principal
type DogRef struct {
	ID  string `header:"id" json:"id"`
	Dog Dog    `header:"inline" json:"dog"`
}

// Dog data model
type Dog struct {
	Name     string `header:"name" json:"name"`
	Breed    string `header:"breed" json:"breed"`
	Color    string `header:"color" json:"color"`
	Birthday string `header:"birthday" json:"birthday"`
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
		// TODO implement type validation and return errors
		return nil
	case "es.hotdoggi.events.dog_updated":
		// TODO implement type validation and return errors
		return nil
	default:
		return fmt.Errorf("unrecognized type name received: %s", typeName)
	}
}
