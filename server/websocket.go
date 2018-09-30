package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{}

func wsserver(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Error("read: %s", err.Error())
			return
		}

		// Print the message to the console
		fmt.Printf("%s says: %s\n", conn.RemoteAddr(), string(msg))

		// Write message back
		if err = conn.WriteMessage(msgType, []byte("Thanks")); err != nil {
			log.Error("write %s", err.Error())
			return
		}
	}
}
