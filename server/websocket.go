package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// var clients = make(map[*websocket.Conn]bool)
// var upgrader = websocket.Upgrader{}

// func wsserver(w http.ResponseWriter, r *http.Request) {
// 	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

// 	for {
// 		// Read message from clients
// 		msgType, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Error("read: %s", err.Error())
// 			return
// 		}

// 		if msg != nil {

// 			// Print the message to the console
// 			fmt.Printf("%s says: %s\n", conn.RemoteAddr(), string(msg))

// 			// Write message back
// 			if err = conn.WriteMessage(msgType, []byte("Thanks")); err != nil {
// 				log.Error("write %s", err.Error())
// 				return
// 			}
// 		}
// 	}
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Directory to watch
	dir string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		// Read message from worker
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Errorln("read error: %s", err.Error())
			return
		}

		if message != nil {

			// Print the message to the console
			fmt.Printf("%s says: %s\n", c.conn.RemoteAddr(), string(message))
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func wsserver(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// func callClients(w http.ResponseWriter, r *http.Request) {
// conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

// if err = conn.WriteMessage(msgType, []byte("Thanks")); err != nil {
// 	log.Error("write %s", err.Error())
// 	return
// }

// for {
// 	// Read message from browser
// 	msgType, msg, err := conn.ReadMessage()
// 	if err != nil {
// 		log.Error("read: %s", err.Error())
// 		return
// 	}

// 	if msg != nil {

// 		// Print the message to the console
// 		fmt.Printf("%s says: %s\n", conn.RemoteAddr(), string(msg))

// 		// Write message back
// 		if err = conn.WriteMessage(msgType, []byte("Thanks")); err != nil {
// 			log.Error("write %s", err.Error())
// 			return
// 		}
// 	}
// }
// }
