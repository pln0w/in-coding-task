package main

import (
	"fmt"
	"log"
	"net/http"
)

const defaultPort = 8080

func main() {

	// Create HTTP server router
	r := NewRouter()

	// Listen for HTTP requests
	log.Printf("listening on port %d", defaultPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), r); err != nil {
		log.Printf("error: %v\n", err.Error())
	}
}
