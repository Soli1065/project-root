package api

import (
	"net/http"
	"project-root/internal/category"
	"project-root/internal/content"
	"project-root/internal/recommendation"
	"project-root/internal/user"

	"github.com/gorilla/mux"
)

func GatewayHandler() http.Handler {
	r := mux.NewRouter()

	// User service routes
	r.HandleFunc("/users/{id}", user.GetUserHandler).Methods("GET")

	// Content service routes
	r.HandleFunc("/content/{id}", content.GetContentHandler).Methods("GET")

	// Category service routes
	r.HandleFunc("/categories/{id}", category.GetCategoryHandler).Methods("GET")

	// Recommendation service routes
	r.HandleFunc("/recommendations/{userId}", recommendation.GetRecommendationHandler).Methods("GET")

	return r
}
