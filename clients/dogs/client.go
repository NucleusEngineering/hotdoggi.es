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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/sacOO7/gowebsocket"

	dogs "github.com/helloworlddan/hotdoggi.es/lib/dogs"
)

// Websockets endpoint
const endpoint = "wss://api.hotdoggies.stamer.demo.altostrat.com/v1/dogs/"

const (
	colorYellow = "\033[93m"
	colorRed    = "\033[91m"
	colorWhite  = "\033[0m"
)

func main() {
	// Switch to subscribe for single dog updates only
	var dogID = flag.String("d", "", "get message for a single dog")
	flag.Parse()

	// Catch CTRL-C signals to exit streaming
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New(fmt.Sprintf("%s%s", endpoint, *dogID))
	socket.ConnectionOptions = gowebsocket.ConnectionOptions{
		UseCompression: false,
		Subprotocols:   []string{"chat", "superchat"},
	}

	// JWT access token for API access
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalln("no credentials provided in $TOKEN")
	}
	socket.RequestHeader.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	socket.RequestHeader.Set("Accept-Encoding", "gzip, deflate, sdch")
	socket.RequestHeader.Set("Accept-Language", "en-US,en;q=0.8")
	socket.RequestHeader.Set("Pragma", "no-cache")
	socket.RequestHeader.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.87 Safari/537.36")

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Fatal(">>> Received connect error: ", err)
	}
	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println(">>> Connected to server")
	}
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {

		// Deserialize received JSON update
		var ref dogs.DogRef
		err := json.Unmarshal([]byte(message), &ref)
		if err != nil {
			log.Println(">>> Failed to deserialize dog update")
		}

		// Pretty print received update to STDOUT
		fmt.Printf("%s\t%s%s%s\tmoved to %s(%f,%f)%s\n", ref.ID, colorRed, ref.Dog.Name, colorWhite, colorYellow, ref.Dog.Location.Latitude, ref.Dog.Location.Longitude, colorWhite)
	}
	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println(">>> Disconnected from server")
	}
	socket.Connect()

	// Listener for CTRL-C exit signals
	for range interrupt {
		fmt.Println()
		log.Println(">>> Exiting ...")
		socket.Close()
		return
	}
}
