package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// Route structure
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes - slice of routes
type Routes []Route

var (
	handler = NewHandlerController()

	routes = Routes{
		Route{
			"Index", "GET",
			"/", handler.Index,
		},
		Route{
			"HealthCheck", "GET",
			"/health", handler.HealthCheck,
		},
	}
)

// LogRequest - logs each request details
func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// NewRouter - creates new Mux Router instance
// and registers handlers and middleware for each route
func NewRouter() *mux.Router {

	// Create new router object
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {

		// Register route in router
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(LogRequest(route.HandlerFunc))
	}

	router.HandleFunc("/ws/test", wsserver)

	return router
}
