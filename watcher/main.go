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
	log.SetFormatter(&log.TextFormatter{})

	// Set log output
	file, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.SetOutput(os.Stdout)
		log.Info("Failed to log to file, using default stderr")
	}

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

	// Prepare WebSocket address
	var wsaddr = flag.String("wsaddr", fmt.Sprintf("%s:%s", lbhost, lbport), "WebSocker service URL")
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *wsaddr, Path: "/ws/test"}
	fmt.Printf("WATCHER %s connecting to [websocket: %s]", hostname, u.String())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error(err.Error())
	}
	defer c.Close()

	fmt.Printf("REMOTE ADDR: #%v", c.RemoteAddr)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Error("read %s", err.Error())
				return
			}
			fmt.Printf("revices: %s\n", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Error("write %s", err.Error())
				return
			}
		case <-interrupt:
			log.Error("WebSocket interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error("write close %s", err.Error())
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
