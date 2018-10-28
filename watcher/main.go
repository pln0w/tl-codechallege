package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

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

	dir := os.Getenv("DIR_PATH")
	if os.Getenv("DIR") != "" {
		dir = os.Getenv("DIR")
	}

	// Prepare WebSocket connection URL
	var wsaddr = flag.String("wsaddr", fmt.Sprintf("%s:%s", lbhost, lbport), "WebSocker service URL")
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *wsaddr, Path: "/ws/test"}
	hostname, _ := os.Hostname()

	fmt.Printf("WATCHER %s is connecting to %s\n", hostname, u.String())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Dial websockets server
	var header http.Header
	header = make(http.Header)
	header.Add("dir", dir)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		fmt.Printf("[dial error]: %v\n", err.Error())
	}
	// defer conn.Close()

	done := make(chan struct{})
	respmsg := make(chan []byte)

	// Concurrently run reading messages
	go func() {
		defer close(done)

		for {
			// Read messages forever
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("[read error]: %v\n", err.Error())
				return
			}

			if message != nil {
				// Pass message via channel, if any
				respmsg <- message
				fmt.Printf("reviced: %s\n", message)
			}
		}
	}()

	// For ever listen on channels
	for {
		select {
		case <-done:
			return
		case t := <-respmsg:
			// Send text data found at channel
			bytes := []byte(t)
			err := conn.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				fmt.Printf("[write error]: %s\n", err.Error())
				return
			}
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			fmt.Printf("[websocket interrupt error]\n")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Printf("[write error]: %v\n", err.Error())
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
