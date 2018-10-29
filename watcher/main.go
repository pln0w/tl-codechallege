package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var fileNames []string

func main() {

	fileNames = []string{}

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

	u := url.URL{Scheme: "ws", Host: *wsaddr, Path: "/ws/register"}
	hostname, _ := os.Hostname()

	fmt.Printf("watcher %s for dir [%s]\n", hostname, dir)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Dial websockets server
	var header http.Header
	header = make(http.Header)
	header.Add("dir", dir)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		fmt.Printf("[ERROR] (dial error): %v\n", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("host %s connected\n", u.String())

	done := make(chan struct{})
	cmd := make(chan []byte)

	// Concurrently run reading messages
	go func() {
		defer close(done)

		for {
			// Read messages from server
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("[ERROR] (read error): %v\n", err.Error())
				return
			}

			if message != nil {
				// Pass message via channel, if any
				cmd <- message
			}
		}
	}()

	// For ever listen on channels
	for {
		select {
		case <-done:
			return
		case command := <-cmd:
			if string(command) == dir {
				files, err := ioutil.ReadDir(dir)
				if err != nil {
					fmt.Printf("[ERROR] (reading directory error): %v\n", err.Error())
					os.Exit(-1)
				}
				for _, f := range files {
					fileNames = append(fileNames, f.Name())
				}
				wErr := conn.WriteJSON(fileNames)
				if wErr != nil {
					fmt.Printf("[ERROR] (reply error): %v\n", wErr.Error())
					return
				}
				fileNames = []string{}
			}
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			fmt.Printf("[ERROR] (websocket interrupt error)\n")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Printf("[ERROR] (write error): %v\n", err.Error())
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
