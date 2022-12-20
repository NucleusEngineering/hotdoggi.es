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
	"net/http"

	gin "github.com/gin-gonic/gin"
	websocket "github.com/gorilla/websocket"
	trace "go.opentelemetry.io/otel/trace"

	dogs "github.com/helloworlddan/hotdoggi.es/lib/dogs"
)

var upgrader = websocket.Upgrader{
	// Allow all origins
	CheckOrigin: func(r *http.Request) bool { return true },
}

// OptionsHandler for CORS preflights OPTIONS requests
func OptionsHandler(c *gin.Context) {}

// ListHandler implements GET /
func ListHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.handler:list")
	defer span.End()

	user := c.MustGet("principal").(*dogs.Principal).ID
	streaming := false
	if upgradeHeader, ok := c.Request.Header["Upgrade"]; ok {
		for _, upgrade := range upgradeHeader {
			if upgrade == "websocket" {
				streaming = true
			}
		}
	}

	if !streaming {
		result, err := List(ctx, user)
		if err != nil {
			log.Printf("error: %v\n", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": fmt.Sprintf("failed to retrieve objects: %v", err),
			})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}

	// Upgrade to websocket stream
	ctx, span = (*tracer).Start(ctx, "dogs.handler:list#streaming")
	defer span.End()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	dogs := make(chan dogs.DogRef)
	go ListStream(ctx, user, dogs)

	for {
		//Response message to client
		dogRef, ok := <-dogs
		if !ok {
			conn.WriteJSON(fmt.Errorf("failed to stream dogs"))
			conn.Close()
		}
		err = conn.WriteJSON(dogRef)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

// GetHandler implements GET /{key}
func GetHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.handler:get")
	defer span.End()

	user := c.MustGet("principal").(*dogs.Principal).ID

	key := c.Param("key")
	streaming := false
	if upgradeHeader, ok := c.Request.Header["Upgrade"]; ok {
		for _, upgrade := range upgradeHeader {
			if upgrade == "websocket" {
				streaming = true
			}
		}
	}

	if !streaming {
		result, err := Get(ctx, user, key)
		if err != nil {
			log.Printf("error: %v\n", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": fmt.Sprintf("failed to retrieve object: %v", err),
			})
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}

	// Upgrade to websocket stream
	ctx, span = (*tracer).Start(ctx, "dogs.handler:get#streaming")
	defer span.End()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	dogs := make(chan dogs.DogRef)
	go GetStream(ctx, user, key, dogs)

	for {
		//Response message to client
		dogRef, ok := <-dogs
		if !ok {
			conn.WriteJSON(fmt.Errorf("failed to stream updates for dog"))
			conn.Close()
		}
		err = conn.WriteJSON(dogRef)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

// EventHandler implements POST /
func EventHandler(c *gin.Context) {
	ctx := c.MustGet("trace.context").(context.Context)
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.handler:event")
	defer span.End()

	typeName := c.MustGet("event.type").(string)
	caller := c.MustGet("principal").(*dogs.Principal)
	ref := c.MustGet("event.data").(*dogs.DogRef)
	switch typeName {
	case "es.hotdoggi.events.dog_added":
		err := dogAdded(ctx, c, caller, ref.Dog)
		if err != nil {
			log.Printf("error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to add dog: %v", err),
			})
		}
		c.Status(http.StatusCreated)
	case "es.hotdoggi.events.dog_removed":
		err := dogRemoved(ctx, c, caller, ref)
		if err != nil {
			log.Printf("error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to remove dog: %v", err),
			})
		}
		c.Status(http.StatusAccepted)
	case "es.hotdoggi.events.dog_updated":
		err := dogUpdated(ctx, c, caller, ref)
		if err != nil {
			log.Printf("error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to update dog: %v", err),
			})
		}
		c.Status(http.StatusAccepted)
	case "es.hotdoggi.events.dog_moved":
		err := dogMoved(ctx, c, caller, ref, ref.Dog.Location.Latitude, ref.Dog.Location.Longitude)
		if err != nil {
			log.Printf("error: %v\n", err)
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

func dogAdded(ctx context.Context, c *gin.Context, caller *dogs.Principal, dog dogs.Dog) error {
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.handler:event.added")
	defer span.End()

	dog.Metadata.Owner = caller.ID

	_, err := Add(ctx, dog)
	return err
}

func dogRemoved(ctx context.Context, c *gin.Context, caller *dogs.Principal, ref *dogs.DogRef) error {
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.handler:event.removed")
	defer span.End()
	_, err := Get(ctx, caller.ID, ref.ID)
	if err != nil {
		return err
	}

	err = Delete(ctx, ref.ID)
	return err
}

func dogUpdated(ctx context.Context, c *gin.Context, caller *dogs.Principal, ref *dogs.DogRef) error {
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.handler:event.updated")
	defer span.End()
	_, err := Get(ctx, caller.ID, ref.ID)
	if err != nil {
		return err
	}
	_, err = Update(ctx, ref.ID, ref.Dog)
	return err
}

func dogMoved(ctx context.Context, c *gin.Context, caller *dogs.Principal, ref *dogs.DogRef, lat float32, long float32) error {
	tracer := Global["client.trace.tracer"].(*trace.Tracer)
	ctx, span := (*tracer).Start(ctx, "dogs.handler:event.moved")
	defer span.End()
	existing, err := Get(ctx, caller.ID, ref.ID)
	if err != nil {
		return err
	}

	existing.Dog.Location.Latitude = lat
	existing.Dog.Location.Longitude = long

	_, err = Update(ctx, ref.ID, existing.Dog)
	return err
}
