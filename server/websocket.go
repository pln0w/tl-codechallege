package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 2
)

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

	// Current list of files
	files []string
}

// function that keeps alive client by periodical ping
func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					fmt.Printf("[ERROR] (write close message error): %v\n", err.Error())
				}
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Printf("[ERROR] (write error): %v\n", err.Error())
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				if websocket.IsUnexpectedCloseError(err) || websocket.IsCloseError(err) {
					c.conn.Close()
				} else {
					fmt.Printf("[ERROR] (watcher unexpectedly crashed): %v\n", err.Error())
				}
				hub.unregister <- c
				return
			}
		}
	}
}

// read function gets messages from the websocket connection
func (c *Client) read() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	var received []string

	for {
		err := c.conn.ReadJSON(&received)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				fmt.Printf("[ERROR] (read error): %v\n", err.Error())
			}
			break
		}

		if len(received) > 0 {
			c.files = received
		}
	}
}

// serveWs handles websocket requests from the peer.
func wsserver(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[ERROR] (upgrader error for client %s): %v", r.RemoteAddr, err.Error())
		return
	}
	defer conn.Close()

	fmt.Printf("new connection accepted from watcher %s\n", r.RemoteAddr)

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.dir = r.Header.Get("dir")

	hub.register <- client

	go client.write()
	client.read()
}
