package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	port := "4000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	// Create router
	router := NewRouter()

	log.Printf("SLAVE listening on %s", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Println("ListenAndServer Error", err)
	}
}
