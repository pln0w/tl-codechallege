package main

import (
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"

	"net/http"
	"os"
)

var hub *Hub

func init() {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true})

	// Log to docker container output
	log.SetOutput(os.Stdout)

	// Set info log level
	log.SetLevel(log.InfoLevel)
}

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
	if err := http.ListenAndServe(*addr, router); err != nil {
		log.Error(err.Error())
	}

	fmt.Printf("SERVER %s listening on [http: %s]\n", os.Getenv("WHOAMI"), port)
}
