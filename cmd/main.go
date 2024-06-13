// cmd/main.go

package main

import (
	"fmt"
	// "log"
	"net/http"

	"project-root/internal/api_gateway"

	"project-root/internal/config"
	"project-root/internal/pkg/middleware"

	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
}

func main() {
	log.Info("Starting the server...")
	fmt.Printf("Server is starting... \n")

	// Initialize the configuration and database connection
	fmt.Printf("Initializing Database... \n")

	config.Initialize()

	// Create a new router
	router := mux.NewRouter()

	// Serve static files for HLS
	router.PathPrefix("/hls/").Handler(http.StripPrefix("/hls/", http.FileServer(http.Dir("./hls"))))

	// Set up API routes and pass the database instance
	api_gateway.SetAPIRoutes(router, config.DB)

	// Wrap the router with the CORS middleware
	corsRouter := middleware.CORSMiddleware(router)

	// // Initialize a new Gorilla mux router
	// router := api_gateway.NewRouter()

	// // Set up API routes
	// api_gateway.SetAPIRoutes(router)

	// Run the HTTP server
	addr := ":8080"
	fmt.Printf("Server is running on %s\n", addr)
	if err := api_gateway.RunServer(corsRouter, addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
