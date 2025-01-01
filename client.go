package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// / Represents a single chatting user
type client struct {
	/// ws for this client
	socket *websocket.Conn
	/// channel to receive messages from other clients
	receive chan []byte
	/// room this client is chatting in
	room *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		msgBytes := []byte(msg)
		c.room.forward <- &msgBytes
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			fmt.Printf("Error %s occured while writing message", err)
		}
	}
}
