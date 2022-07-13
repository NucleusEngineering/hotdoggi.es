package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/sacOO7/gowebsocket"
)

const endpoint = "wss://api.hotdoggies.stamer.demo.altostrat.com/dogs/"

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

const (
	colorYellow = "\033[93m"
	colorRed    = "\033[91m"
	colorWhite  = "\033[0m"
)

func main() {
	var dogID = flag.String("d", "", "get message for a single dog")
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New(fmt.Sprintf("%s%s", endpoint, *dogID))
	socket.ConnectionOptions = gowebsocket.ConnectionOptions{
		UseCompression: false,
		Subprotocols:   []string{"chat", "superchat"},
	}

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

		var ref DogRef
		err := json.Unmarshal([]byte(message), &ref)
		if err != nil {
			log.Println(">>> Failed to deserialize dog update")
		}

		fmt.Printf("%s\t%s%s%s\tmoved to %s(%f,%f)%s\n", ref.ID, colorRed, ref.Dog.Name, colorWhite, colorYellow, ref.Dog.Location.Latitude, ref.Dog.Location.Longitude, colorWhite)
	}
	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println(">>> Disconnected from server")
	}
	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println(">>> Exiting ...")
			socket.Close()
			return
		}
	}
}
