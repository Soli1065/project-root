// internal/content/handler.go
package content

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// GetAllContentsHandler handles the request to retrieve all contents
func GetAllContentsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contents, err := GetAllContents(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(contents)
	}
}

// GetContentByIDHandler handles the request to retrieve a content by its ID
func GetContentByIDHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		content, err := GetContentByID(db, uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(content)
	}
}

// CreateContentHandler handles the request to create a new content
func CreateContentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var content Content
		if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := CreateContent(db, &content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(content)
	}
}

// GetAllContentHandler retrieves all content items from the database
func GetAllContentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var content []Content
		if err := db.Find(&content).Error; err != nil {
			http.Error(w, "Failed to fetch content", http.StatusInternalServerError)
			return
		}

		// Respond with the content data
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(content)
	}
}

// UpdateContentHandler handles the request to update an existing content
func UpdateContentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var content Content
		if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		content.ID = uint(id)
		if err := UpdateContent(db, &content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(content)
	}
}

// DeleteContentHandler handles the request to delete a content
func DeleteContentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var content Content
		if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		content.ID = uint(id)

		if err := DeleteContent(db, &content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

type VideoUploadResponse struct {
	ID       uint   `json:"id"`
	VideoURL string `json:"video_url"`
}

func UploadVideoHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form values
		title := r.FormValue("title")
		description := r.FormValue("description")
		authorID, err := strconv.Atoi(r.FormValue("author_id"))
		if err != nil {
			http.Error(w, "Invalid author ID", http.StatusBadRequest)
			return
		}
		authorName := r.FormValue("author_name")
		categoryID, err := strconv.Atoi(r.FormValue("category_id"))
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		// Handle file upload
		file, handler, err := r.FormFile("video")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save the uploaded video file temporarily
		tempFilePath := fmt.Sprintf("/tmp/%d%s", time.Now().Unix(), filepath.Ext(handler.Filename))
		tempFile, err := os.Create(tempFilePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()
		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert video to HLS format
		hlsOutputPath := fmt.Sprintf("/var/www/html/hls/videos/%d", time.Now().Unix())
		os.MkdirAll(hlsOutputPath, os.ModePerm)
		ffmpegCmd := exec.Command("ffmpeg", "-i", tempFilePath, "-codec: copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", filepath.Join(hlsOutputPath, "index.m3u8"))
		if err := ffmpegCmd.Run(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Store content record in database
		videoURL := fmt.Sprintf("/hls/%d/index.m3u8", time.Now().Unix())
		content := ContentModel{
			Title:       title,
			Description: description,
			URL:         videoURL,
			CategoryID:  uint(categoryID),
			AuthorID:    uint(authorID),
			AuthorName:  authorName,
			CreatedAt:   time.Now(),
		}
		if err := db.Create(&content).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the video URL
		response := VideoUploadResponse{
			ID:       content.ID,
			VideoURL: videoURL,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
