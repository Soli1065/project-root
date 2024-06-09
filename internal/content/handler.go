// internal/content/handler.go
package content

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		// Set a maximum upload size of 100 MB
		r.Body = http.MaxBytesReader(w, r.Body, 100<<20) // 100MB

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
			http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Check file size
		if handler.Size > 100<<20 { // 100MB
			http.Error(w, "File size exceeds 100MB limit", http.StatusBadRequest)
			return
		}

		// Save the uploaded video file temporarily
		tempFilePath := fmt.Sprintf("/tmp/%d%s", time.Now().Unix(), filepath.Ext(handler.Filename))
		tempFile, err := os.Create(tempFilePath)
		if err != nil {
			http.Error(w, "Could not create temp file", http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()
		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, "Could not copy file to temp location", http.StatusInternalServerError)
			return
		}

		// Handle image file upload
		imageFile, imageHeader, err := r.FormFile("image")
		var imageURL string
		if err == nil {
			defer imageFile.Close()

			// Save the uploaded image file
			imageFilePath := fmt.Sprintf("/var/www/academyserverapp/images/%d%s", time.Now().Unix(), filepath.Ext(imageHeader.Filename))
			imageFileDest, err := os.Create(imageFilePath)
			if err != nil {
				http.Error(w, "Error creating image file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer imageFileDest.Close()
			_, err = io.Copy(imageFileDest, imageFile)
			if err != nil {
				http.Error(w, "Error saving the image file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			imageURL = imageFilePath
		} else if err != http.ErrMissingFile {
			http.Error(w, "Error retrieving the image file: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Get video duration using ffmpeg
		output, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", tempFilePath).Output()
		if err != nil {
			http.Error(w, "Error getting video duration: "+err.Error(), http.StatusInternalServerError)
			return
		}
		duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
		if err != nil {
			http.Error(w, "Error parsing video duration: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert video to HLS format
		hlsOutputPath := fmt.Sprintf("/var/www/academyserverapp/hls/videos/%d", time.Now().Unix())
		err = os.MkdirAll(hlsOutputPath, os.ModePerm)
		if err != nil {
			http.Error(w, "Could not create HLS output directory", http.StatusInternalServerError)
			return
		}

		ffmpegCmd := exec.Command("ffmpeg", "-i", tempFilePath, "-codec:", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", filepath.Join(hlsOutputPath, "index.m3u8"))
		ffmpegOutput, err := ffmpegCmd.CombinedOutput()
		if err != nil {
			log.Printf("FFmpeg command failed with error: %s\nOutput: %s", err, string(ffmpegOutput))
			http.Error(w, "Could not process video", http.StatusInternalServerError)
			return
		}

		// Store content record in database
		videoURL := fmt.Sprintf("/hls/videos/%d/index.m3u8", time.Now().Unix())
		content := Content{
			Title:       title,
			Description: description,
			URL:         videoURL,
			CategoryID:  uint(categoryID),
			AuthorID:    uint(authorID),
			AuthorName:  authorName,
			CreatedAt:   time.Now(),
			ImageURL:    imageURL,
			Duration:    uint(duration),
			IsLive:      false,
			ViewCount:   0,
		}
		if err := db.Create(&content).Error; err != nil {
			http.Error(w, "Could not save content record", http.StatusInternalServerError)
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

func IncrementViewCountHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentID, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "Invalid content ID", http.StatusBadRequest)
			return
		}

		var content Content
		if err := db.First(&content, contentID).Error; err != nil {
			http.Error(w, "Content not found", http.StatusNotFound)
			return
		}

		content.ViewCount++
		if err := db.Save(&content).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
