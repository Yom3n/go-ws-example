package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// holds all current clients in this room
	clients map[*client]struct{}
	// channel for clients wishing to join the room
	join chan *client
	// channel for clients whising to leave the room
	leave chan *client
	// channel that holds incoming messages that should be forwarded to other clients
	forward chan *[]byte
}

func NewRoom() *room {
	return &room{
		clients: make(map[*client]struct{}),
		join:    make(chan *client),
		leave:   make(chan *client),
		forward: make(chan *[]byte),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = struct{}{}
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		case msg := <-r.forward:
			for c := range r.clients {
				c.receive <- *msg
			}

		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(w, request, nil)
	if err != nil {
		log.Fatalf("Error occured while upgrading to websocket: %s", err)
		return
	}
	defer func() {
		log.Println("Closing websocket connection with the room")
		conn.Close()
	}()

	client := &client{
		socket:  conn,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
	}
	r.join <- client

	defer func() {
		r.leave <- client
	}()

	go client.write()
	client.read()

}
