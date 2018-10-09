package main

import (
	"flag"
	"fmt"
	"net/url"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"os"
)

func init() {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})

	// Log to docker container output
	log.SetOutput(os.Stdout)

	// Set info log level
	log.SetLevel(log.InfoLevel)
}

func main() {

	// Define WebSocket server port
	lbport := "80"
	if os.Getenv("LB_PORT") != "" {
		lbport = os.Getenv("LB_PORT")
	}

	// Define WebSocket server host
	lbhost := ""
	if os.Getenv("LB_HOST") != "" {
		lbhost = os.Getenv("LB_HOST")
	}

	hostname, _ := os.Hostname()

	// Prepare WebSocket connection URL
	var wsaddr = flag.String("wsaddr", fmt.Sprintf("%s:%s", lbhost, lbport), "WebSocker service URL")
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *wsaddr, Path: "/ws/test"}
	fmt.Printf("WATCHER %s connecting to [websocket: %s]\n", hostname, u.String())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Dial websockets server
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error(err.Error())
	}
	defer c.Close()

	done := make(chan struct{})
	respmsg := make(chan []byte)

	// Concurrently run reading messages
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Errorln("read error: %s", err.Error())
				return
			}

			if message != nil {
				respmsg <- message
				fmt.Printf("revices: %s\n", message)
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case t := <-respmsg:
			bytes := []byte(t)

			err := c.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				log.Errorln("write error: %s", err.Error())
				return
			}
		case <-interrupt:
			log.Println("websocket interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Errorln("write error: %s", err.Error())
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
