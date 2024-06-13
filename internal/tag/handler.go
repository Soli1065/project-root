package tag

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	// "os"
)

func CreateTagHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tag Tag
		if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := db.Create(&tag).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tag)
	}
}

func GetTagsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tags []Tag
		if err := db.Find(&tags).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tags)
	}
}

func UpdateTagHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tag Tag
		id := mux.Vars(r)["id"]
		if err := db.First(&tag, id).Error; err != nil {
			http.Error(w, "Tag not found", http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := db.Save(&tag).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tag)
	}
}

func DeleteTagHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if err := db.Delete(&Tag{}, id).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
