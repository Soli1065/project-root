// internal/recommendation/handler.go
package recommendation

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetAllRecommendationsHandler handles the request to retrieve all recommendations
func GetAllRecommendationsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recommendations, err := GetAllRecommendations(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(recommendations)
	}
}

// GetRecommendationByIDHandler handles the request to retrieve a recommendation by its ID
func GetRecommendationByIDHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		recommendation, err := GetRecommendationByID(db, uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(recommendation)
	}
}

// CreateRecommendationHandler handles the request to create a new recommendation
func CreateRecommendationHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var recommendation Recommendation
		if err := json.NewDecoder(r.Body).Decode(&recommendation); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := CreateRecommendation(db, &recommendation); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(recommendation)
	}
}

// UpdateRecommendationHandler handles the request to update an existing recommendation
func UpdateRecommendationHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var recommendation Recommendation
		if err := json.NewDecoder(r.Body).Decode(&recommendation); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		recommendation.ID = uint(id)
		if err := UpdateRecommendation(db, &recommendation); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(recommendation)
	}
}

// DeleteRecommendationHandler handles the request to delete a recommendation
func DeleteRecommendationHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := DeleteRecommendation(db, uint(id)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
