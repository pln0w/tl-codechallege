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

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
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
}

// read function gets messages from the websocket connection
func (c *Client) read() {
	defer func() {
		c.conn.Close()
	}()

	for {
		_, received, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				fmt.Printf("[read error]: %v\n", err.Error())
			}
			break
		}

		fmt.Printf("received: %s\n", received)

		// if received == "pattern" { do() }
	}
}

// write function sends messages to the websocket connection.
func (c *Client) write(hub *Hub) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					fmt.Printf("[write error]: %v\n", err.Error())
					break
				}
				break
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Printf("[next writer error]: %v\n", err.Error())
				break
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				fmt.Printf("[writer close error]: %v\n", err.Error())
				break
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				if websocket.IsUnexpectedCloseError(err) || websocket.IsCloseError(err) {
					c.conn.Close()
				} else {
					fmt.Printf("[watcher crash]: %v\n", err.Error())
				}
				hub.unregister <- c
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func wsserver(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[client %s error]: %v", r.RemoteAddr, err.Error())
		return
	}
	defer conn.Close()

	fmt.Printf("new connection accepted from watcher %s\n", r.RemoteAddr)

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.dir = r.Header.Get("dir")

	hub.register <- client

	go client.write(hub)
	client.read()

	// _, _, readErr := client.conn.ReadMessage()
	// if readErr != nil {
	// 	fmt.Printf("[watcher crash]: %v\n", readErr.Error())
	// }
}
