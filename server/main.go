package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
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
	fmt.Printf("SERVER %s is listening at port %s\n", hostname, port)

	if err := http.ListenAndServe(*addr, router); err != nil {
		fmt.Printf("[server error]: %v\n", err.Error())
	}
}
