package main

import "github.com/gorilla/websocket"

/// Represents a single chatting user
type client struct {
	/// ws for this client
	socket *websocket.Conn
	/// channel to receive messages from other clients
	receive chan []byte
	/// room this client is chatting in
	room *room
}

func (c *client) read(){
	defer c.socket.Close()
	for {
	type, msg, err := c.socket.ReadMessage()
	if err != nil {
		return
	}
	}
}