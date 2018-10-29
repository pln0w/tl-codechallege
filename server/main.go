package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

var hub *Hub

func main() {

	// Define HTTP server port
	port := "80"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	// Prepare HTTP server address
	var addr = flag.String("addr", fmt.Sprintf("0.0.0.0:%s", port), "HTTP service URL")
	flag.Parse()

	hub = newHub()
	go hub.run()

	// Create router
	router := NewRouter()

	// Serve over HTTP
	hostname, _ := os.Hostname()
	fmt.Printf("server %s is listening at port %s\n", hostname, port)

	// Concurrently run watchers updating their files lists
	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			for c := range hub.clients {
				go func() {
					c.send <- []byte(c.dir)
				}()
			}
		}
	}()

	// Listen for new clients and HTTP requests
	if err := http.ListenAndServe(*addr, router); err != nil {
		fmt.Printf("[ERROR] (server error): %v\n", err.Error())
	}
}
