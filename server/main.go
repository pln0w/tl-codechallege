package main

import (
	"flag"
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
	if err != nil {
		log.SetOutput(os.Stdout)
		log.Info("Failed to log to file, using default stderr")
	} else {
		log.SetOutput(file)
	}

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

	// Define WebSocket server port
	// wsport := "12345"
	// if os.Getenv("WS_PORT") != "" {
	// 	wsport = os.Getenv("WS_PORT")
	// }

	// // Prepare WebSocket address
	// var wsaddr = flag.String("wsaddr", fmt.Sprintf("0.0.0.0:%s", wsport), "WebSocker service URL")

	flag.Parse()

	// Create router
	router := NewRouter()

	// Serve over HTTP
	if err := http.ListenAndServe(*addr, router); err != nil {
		log.Error(err.Error())
	}

	// // Register WS handler
	// http.HandleFunc("/ws/test", wsserver)

	// // Serve over WebSocket (internal communication)
	// if err := http.ListenAndServe(*wsaddr, nil); err != nil {
	// 	log.Error(err.Error())
	// }

	log.Printf("SERVER %s listening on [http: %s]", os.Getenv("WHOAMI"), port)
}
