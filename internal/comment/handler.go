package comment

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CreateCommentHandler handles the creation of a new comment
func CreateCommentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newComment Comment
		if err := json.NewDecoder(r.Body).Decode(&newComment); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		contentID, err := strconv.ParseUint(mux.Vars(r)["contentID"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newComment.ContentID = uint(contentID)

		if err := db.Create(&newComment).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newComment)
	}
}

// GetCommentsHandler handles retrieving all comments for a content
func GetCommentsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var comments []Comment
		contentID, err := strconv.ParseUint(mux.Vars(r)["contentID"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := db.Where("content_id = ?", uint(contentID)).Find(&comments).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}

// UpdateCommentHandler handles updating an existing comment
func UpdateCommentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updatedComment Comment
		if err := json.NewDecoder(r.Body).Decode(&updatedComment); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		commentID, err := strconv.ParseUint(mux.Vars(r)["commentID"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := db.Model(&Comment{}).Where("id = ?", uint(commentID)).Updates(updatedComment).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedComment)
	}
}

// DeleteCommentHandler handles deleting a comment
func DeleteCommentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commentID, err := strconv.ParseUint(mux.Vars(r)["commentID"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := db.Delete(&Comment{}, uint(commentID)).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
