package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/aarthikrao/jumpAndShoot/objects"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:4000", "http service address")

func main() {
	flag.Parse()

	http.HandleFunc("/gameEngine", handler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "JumpAndShoot engine is running on this port")
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

func handler(w http.ResponseWriter, r *http.Request) {
	playerName := r.URL.Query()["id"][0]
	log.Println("PlayerJoined :", playerName)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ERROR", err)
		return
	}

	// Create a new player instance
	p := objects.NewPlayer(playerName, conn)

	// Start listening to messages from player
	go p.RecieveMessages()

	// Find a match for the player
	objects.GS.Match(p)
}
