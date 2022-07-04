package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/sacOO7/gowebsocket"
)

const endpoint = "wss://api.hotdoggies.stamer.demo.altostrat.com/dogs/"

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
	socket.RequestHeader.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	socket.RequestHeader.Set("Accept-Encoding", "gzip, deflate, sdch")
	socket.RequestHeader.Set("Accept-Language", "en-US,en;q=0.8")
	socket.RequestHeader.Set("Pragma", "no-cache")
	socket.RequestHeader.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.87 Safari/537.36")

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Fatal("Received connect error ", err)
	}
	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server")
	}
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Received message  " + message)
	}
	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	}
	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server")
	}
	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}
}
