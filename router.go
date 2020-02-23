package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Route structure
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes type - slice of Route objects
type Routes []Route

var (
	apiController = NewAPIController()
	routes        = Routes{
		Route{
			"health check", "GET", "/", apiController.HealthCheck,
		},
		Route{
			"routes", "GET", "/routes", apiController.Routes,
		},
	}
)

// LogRequest - logs each request details
func LogRequest(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("[HTTP] (%s %s %s)\n", r.Method, r.URL, r.RemoteAddr)

		handler.ServeHTTP(w, r)
	})
}

// NewRouter - creates new Mux Router instance. Bind handlers and middleware for each route
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {

		// Bind route
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(LogRequest(route.HandlerFunc))
	}

	return router
}
