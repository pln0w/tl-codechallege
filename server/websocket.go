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

// readPump pumps messages from the websocket connection to the hub.
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	// defer c.conn.Close()
	_, _, err := c.conn.ReadMessage()
	if err != nil {
		return
	}
	// defer func() {
	// 	c.hub.unregister <- c
	// 	c.conn.Close()
	// }()

	// c.conn.SetReadLimit(maxMessageSize)
	// c.conn.SetReadDeadline(time.Now().Add(time.Hour * 2242560))
	// c.conn.SetPongHandler(func(string) error {
	// 	c.conn.SetReadDeadline(time.Now().Add(time.Hour * 2242560))
	// 	return nil
	// })

	// time.Sleep(2 * time.Second)
	// // Read messages from worker forever and print to
	// // standard output with worker IP address
	// for {
	// 	_, message, err := c.conn.ReadMessage()
	// 	if err != nil {
	// 		fmt.Printf("[read error]: %v\n", err.Error())
	// 		break
	// 	}

	// 	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	// 	fmt.Printf("%s says: %s\n", c.conn.RemoteAddr(), string(message))
	// }
}

// writePump pumps messages from the hub to the websocket connection.
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump(hub *Hub) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
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
				fmt.Printf("[close error]: %v\n", err.Error())
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

	fmt.Printf("==> new connection accepted from watcher %s\n", r.RemoteAddr)

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	hub.register <- client

	go client.writePump(hub)

	_, _, readErr := client.conn.ReadMessage()
	if readErr != nil {
		fmt.Printf("[watcher crash]: %v\n", readErr.Error())
	}
}
