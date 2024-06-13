// cmd/main.go

package main

import (
	"fmt"
	// "log"
	"net/http"

	"project-root/internal/api_gateway"

	"project-root/internal/config"
	"project-root/internal/pkg/middleware"

	"encoding/json"
	"project-root/internal/attachment"
	"project-root/internal/content"

	"gorm.io/gorm"

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

// GetAllContentsHandler handles the request to retrieve all contents
func GetAllContentsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contents, err := content.GetAllContents(db)
		if err != nil {
			log.Error("Failed to get all contents: ", err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update each content object before encoding to JSON
		for i := range contents {

			// Populate attachments list
			attachments, err := GetAttachmentsForContent(db, contents[i].ID)
			if err != nil {
				log.Error("Failed to get attachments for content ID ", contents[i].ID, ": ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			contents[i].Attachments = attachments

			log.WithFields(logrus.Fields{
				"content_id":  contents[i].ID,
				"title":       contents[i].Title,
				"attachments": contents[i].Attachments,
			}).Info("Processed content")

			// Log the content object
			fmt.Printf("Content ID %d: %+v", contents[i].ID, contents[i])
		}

		// Encode contents to JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(contents); err != nil {
			log.Error("Failed to encode JSON: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// GetAttachmentsForContent retrieves all attachments for a given content ID
func GetAttachmentsForContent(db *gorm.DB, contentID uint) ([]attachment.Attachment, error) {
	var attachments []attachment.Attachment
	err := db.Where("content_id = ?", contentID).Find(&attachments).Error
	if err != nil {
		log.Error("Error retrieving attachments for content ID ", contentID, ": ", err)
		return nil, err
	}
	log.Infof("Retrieved %d attachments for content ID %d", len(attachments), contentID)

	return attachments, nil
}
