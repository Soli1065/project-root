// cmd/main.go

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/your-module/internal/api_gateway"
)

func main() {
	// Initialize a new Gorilla mux router
	router := api_gateway.NewRouter()

	// Set up API routes
	api_gateway.SetAPIRoutes(router)

	// Run the HTTP server
	addr := ":8080"
	fmt.Printf("Server is running on %s\n", addr)
	if err := api_gateway.RunServer(router, addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
