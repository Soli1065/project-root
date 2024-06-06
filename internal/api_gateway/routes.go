// internal/api_gateway/routes.go

package api_gateway

import (
	"project-root/internal/auth"
	"project-root/internal/category"
	"project-root/internal/content"
	"project-root/internal/recommendation"
	"project-root/internal/user"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// SetAPIRoutes sets up the API routes for the application.
func SetAPIRoutes(router *mux.Router, db *gorm.DB) {
	// Authentication routes
	router.HandleFunc("/auth/login", auth.LoginHandler(db)).Methods("POST")
	router.HandleFunc("/auth/register", auth.RegisterHandler(db)).Methods("POST")

	// User routes
	router.HandleFunc("/users", user.GetAllUsersHandler(db)).Methods("GET")
	router.HandleFunc("/users/{id}", user.GetUserByIDHandler(db)).Methods("GET")
	router.HandleFunc("/users", user.CreateUserHandler(db)).Methods("POST")
	router.HandleFunc("/users/{id}", user.UpdateUserHandler(db)).Methods("PUT")
	router.HandleFunc("/users/{id}", user.DeleteUserHandler(db)).Methods("DELETE")

	// Content routes
	router.HandleFunc("/content", content.GetAllContentHandler(db)).Methods("GET")
	router.HandleFunc("/content/{id}", content.GetContentByIDHandler(db)).Methods("GET")
	router.HandleFunc("/content", content.CreateContentHandler(db)).Methods("POST")
	router.HandleFunc("/content/{id}", content.UpdateContentHandler(db)).Methods("PUT")
	router.HandleFunc("/content/{id}", content.DeleteContentHandler(db)).Methods("DELETE")

	// Category routes
	router.HandleFunc("/categories", category.GetAllCategoriesHandler(db)).Methods("GET")
	router.HandleFunc("/categories/{id}", category.GetCategoryByIDHandler(db)).Methods("GET")
	router.HandleFunc("/categories", category.CreateCategoryHandler(db)).Methods("POST")
	router.HandleFunc("/categories/{id}", category.UpdateCategoryHandler(db)).Methods("PUT")
	router.HandleFunc("/categories/{id}", category.DeleteCategoryHandler(db)).Methods("DELETE")

	// Recommendation routes
	router.HandleFunc("/recommendations", recommendation.GetAllRecommendationsHandler(db)).Methods("GET")
	router.HandleFunc("/recommendations/{id}", recommendation.GetRecommendationByIDHandler(db)).Methods("GET")
	router.HandleFunc("/recommendations", recommendation.CreateRecommendationHandler(db)).Methods("POST")
	router.HandleFunc("/recommendations/{id}", recommendation.UpdateRecommendationHandler(db)).Methods("PUT")
	router.HandleFunc("/recommendations/{id}", recommendation.DeleteRecommendationHandler(db)).Methods("DELETE")
}
