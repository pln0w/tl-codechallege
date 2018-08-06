package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var (
	ctrl = NewController()

	routes = Routes{
		Route{
			"Index", "GET",
			"/", ctrl.Index,
		},
		Route{
			"HealthCheck", "GET",
			"/health", ctrl.HealthCheck,
		},
	}
)

// Router - creates new Mux Router instance
// and registers handlers and middleware for each route
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}
