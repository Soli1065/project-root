// internal/video/handler.go
package video

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func CreateVideoHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var video Video
		if err := json.NewDecoder(r.Body).Decode(&video); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := db.Create(&video).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(video)
	}
}

func GetVideoByIDHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var video Video
		if err := db.First(&video, id).Error; err != nil {
			http.Error(w, "Video not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(video)
	}
}
