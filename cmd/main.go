// cmd/main.go

package main

import (
	"fmt"
	"log"

	"project-root/internal/api_gateway"

	"project-root/internal/config"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Printf("Server is starting... \n")

	// Initialize the configuration and database connection
	fmt.Printf("Initializing Database... \n")
	config.Initialize()

	// Create a new router
	router := mux.NewRouter()

	// Set up API routes and pass the database instance
	api_gateway.SetAPIRoutes(router, config.DB)

	// // Initialize a new Gorilla mux router
	// router := api_gateway.NewRouter()

	// // Set up API routes
	// api_gateway.SetAPIRoutes(router)

	// Run the HTTP server
	addr := ":8080"
	fmt.Printf("Server is running on %s\n", addr)
	if err := api_gateway.RunServer(router, addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
