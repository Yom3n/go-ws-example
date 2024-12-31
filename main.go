package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	ugrader websocket.Upgrader
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.ugrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error %s when upgrading connection to ws", err)
		return
	}
	defer func() {
		log.Println("Closing connection")
		c.Close()
	}()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error %s when reading message from client", err)
			return
		}
		if mt == websocket.BinaryMessage {
			err := c.WriteMessage(websocket.TextMessage, []byte("Server doesn't support binary messages"))
			if err != nil {
				log.Printf("Error %s while sending message to the client", err)
			}
			return
		}
		log.Printf("Received message %s", string(msg))
		if strings.Trim(string(msg), "\n") != "start" {
			err = c.WriteMessage(websocket.TextMessage, []byte("You did not say magic word"))
			if err != nil {
				log.Printf("Error %s while sending message to the client", err)
				return
			}
			continue
		}
		log.Println("Start responding to client")
		i := 1
		for {
			response := fmt.Sprintf("Notification %d", i)
			err := c.WriteMessage(websocket.TextMessage, []byte(response))
			if err != nil {
				log.Printf("Error %s while sending message to the client", err)
				return
			}
			i++
			time.Sleep(2 * time.Second)
		}
	}
}

func main() {
	wsHandler := webSocketHandler{ugrader: websocket.Upgrader{}}
	http.Handle("/", wsHandler)
	addr := "localhost:8080"
	log.Print("Starting server at ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
