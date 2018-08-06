package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"net/http"
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

	// Define HTTP server port
	port := "4000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	// Create router
	router := NewRouter()

	log.Printf("SLAVE listening on %s", port)

	// Serve over HTTP
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), router); err != nil {
		log.Error(err.Error())
	}
}
