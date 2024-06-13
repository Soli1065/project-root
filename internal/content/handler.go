// internal/content/handler.go
package content

import (
	"encoding/json"
	"log"
	"net/http"
	"project-root/internal/attachment"
	"project-root/internal/category"
	"project-root/internal/comment"

	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"fmt"
	"io"

	// "os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"os"
)

// GetAllContentsHandler handles the request to retrieve all contents
func GetAllContentsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contents []Content

		// Fetch all contents with attachments and categories preloaded
		if err := db.Preload("Attachments").Preload("Categories").Find(&contents).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update each content object with comments
		for i := range contents {
			// Populate comments list
			comments, err := GetCommentsForContent(db, contents[i].ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			contents[i].Comments = comments

			// Log the content object (if needed)
			// fmt.Printf("Content ID %d: %+v", contents[i].ID, contents[i])
		}

		// Encode contents to JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(contents); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// GetAttachmentsForContent retrieves all attachments for a given content ID
func GetAttachmentsForContent(db *gorm.DB, contentID uint) ([]attachment.Attachment, error) {
	var attachments []attachment.Attachment
	err := db.Where("content_id = ?", contentID).Find(&attachments).Error
	if err != nil {
		// log.Error("Error retrieving attachments for content ID ", contentID, ": ", err)
		return nil, err
	}
	// log.Infof("Retrieved %d attachments for content ID %d", len(attachments), contentID)

	return attachments, nil
}

// GetCommentsForContent retrieves all comments for a given content ID
func GetCommentsForContent(db *gorm.DB, contentID uint) ([]comment.Comment, error) {
	var comments []comment.Comment
	err := db.Where("content_id = ?", contentID).Find(&comments).Error
	if err != nil {
		log.Printf("Error retrieving comments for content ID %d: %s", contentID, err)
		return nil, err
	}
	log.Printf("Retrieved %d comments for content ID %d", len(comments), contentID)
	return comments, nil
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

// legacy - old - should remove in future
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
	ImageURL string `json:"image_url"`
	Duration string `json:"duration"`
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
			imageURL = fmt.Sprintf("/images/%d%s", time.Now().Unix(), filepath.Ext(imageHeader.Filename))

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
		durationSeconds, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
		if err != nil {
			http.Error(w, "Error parsing video duration: "+err.Error(), http.StatusInternalServerError)
			return
		}
		duration := formatDuration(durationSeconds)

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
			Duration:    duration,
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
			ImageURL: imageURL,
			Duration: content.Duration,
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

// formatDuration converts a float64 duration in seconds to hh:mm:ss format
func formatDuration(duration float64) string {
	hours := int(duration / 3600)
	minutes := int((duration - float64(hours*3600)) / 60)
	seconds := int(duration - float64(hours*3600) - float64(minutes*60))
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

type UploadResponse struct {
	ID          uint                    `json:"id"`
	ContentURL  string                  `json:"content_url"`
	ImageURL    string                  `json:"image_url"`
	Duration    string                  `json:"duration,omitempty"`
	Attachments []attachment.Attachment `json:"attachments,omitempty"`
}

func UploadContentHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set a maximum upload size of 100 MB
		r.Body = http.MaxBytesReader(w, r.Body, 100<<20) // 100MB

		// Parse form values
		title := r.FormValue("title")
		description := r.FormValue("description")
		categoryID, err := strconv.Atoi(r.FormValue("category_id"))
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		AuthorId, err := strconv.Atoi(r.FormValue("author_id"))
		if err != nil {
			http.Error(w, "Invalid author id ID", http.StatusBadRequest)
			return
		}
		AuthorName := r.FormValue("author_name")

		// Handle primary content file upload
		mainFile, mainFileHeader, err := r.FormFile("main_file")
		if err != nil {
			http.Error(w, "Could not get uploaded main file", http.StatusBadRequest)
			return
		}
		defer mainFile.Close()

		// Check main file size
		if mainFileHeader.Size > 100<<20 { // 100MB
			http.Error(w, "File size exceeds 100MB limit", http.StatusBadRequest)
			return
		}

		// Save the uploaded main file
		mainFilePath := fmt.Sprintf("/var/www/academyserverapp/contents/%d%s", time.Now().Unix(), filepath.Ext(mainFileHeader.Filename))
		mainFileDest, err := os.Create(mainFilePath)
		if err != nil {
			http.Error(w, "Could not create main file", http.StatusInternalServerError)
			return
		}
		defer mainFileDest.Close()
		_, err = io.Copy(mainFileDest, mainFile)
		if err != nil {
			http.Error(w, "Could not copy main file to destination", http.StatusInternalServerError)
			return
		}

		var mainFileURL string
		var duration string
		if strings.ToLower(filepath.Ext(mainFileHeader.Filename)) == ".mp4" {
			// Get video duration using ffmpeg
			output, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", mainFilePath).Output()
			if err != nil {
				http.Error(w, "Error getting video duration: "+err.Error(), http.StatusInternalServerError)
				return
			}
			durationSeconds, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
			if err != nil {
				http.Error(w, "Error parsing video duration: "+err.Error(), http.StatusInternalServerError)
				return
			}
			duration = formatDuration(durationSeconds)

			// Convert video to HLS format
			hlsOutputPath := fmt.Sprintf("/var/www/academyserverapp/hls/videos/%d", time.Now().Unix())
			err = os.MkdirAll(hlsOutputPath, os.ModePerm)
			if err != nil {
				http.Error(w, "Could not create HLS output directory", http.StatusInternalServerError)
				return
			}

			ffmpegCmd := exec.Command("ffmpeg", "-i", mainFilePath, "-codec:", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", filepath.Join(hlsOutputPath, "index.m3u8"))
			ffmpegOutput, err := ffmpegCmd.CombinedOutput()
			if err != nil {
				log.Printf("FFmpeg command failed with error: %s\nOutput: %s", err, string(ffmpegOutput))
				http.Error(w, "Could not process video", http.StatusInternalServerError)
				return
			}

			mainFileURL = fmt.Sprintf("/hls/videos/%d/index.m3u8", time.Now().Unix())
		} else {
			mainFileURL = fmt.Sprintf("/contents/%d%s", time.Now().Unix(), filepath.Ext(mainFileHeader.Filename))
		}

		// Handle image file upload (optional)
		var imageURL string
		imageFile, imageHeader, err := r.FormFile("image")
		if err == nil {
			defer imageFile.Close()
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
			imageURL = fmt.Sprintf("/images/%d%s", time.Now().Unix(), filepath.Ext(imageHeader.Filename))
		} else if err != http.ErrMissingFile {
			http.Error(w, "Error retrieving the image file: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Handle attachments upload (optional)
		attachments := []attachment.Attachment{}
		for i := 0; i < 5; i++ { // Assuming a maximum of 5 attachments for example
			fileKey := fmt.Sprintf("attachment_%d", i+1)
			attachmentFile, attachmentHeader, err := r.FormFile(fileKey)
			if err == http.ErrMissingFile {
				continue
			} else if err != nil {
				http.Error(w, "Error retrieving attachment file: "+err.Error(), http.StatusBadRequest)
				return
			}
			defer attachmentFile.Close()

			attachmentFilePath := fmt.Sprintf("/var/www/academyserverapp/attachments/%d%s", time.Now().Unix(), filepath.Ext(attachmentHeader.Filename))
			attachmentFileDest, err := os.Create(attachmentFilePath)
			if err != nil {
				http.Error(w, "Error creating attachment file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer attachmentFileDest.Close()
			_, err = io.Copy(attachmentFileDest, attachmentFile)
			if err != nil {
				http.Error(w, "Error saving the attachment file: "+err.Error(), http.StatusInternalServerError)
				return
			}

			attachmentURL := fmt.Sprintf("/attachments/%d%s", time.Now().Unix(), filepath.Ext(attachmentHeader.Filename))
			attachmentRecord := attachment.Attachment{
				FilePath: attachmentURL,
				FileType: filepath.Ext(attachmentHeader.Filename),
			}

			if strings.ToLower(filepath.Ext(attachmentHeader.Filename)) == ".mp4" {
				// Get video duration using ffmpeg
				// output, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", attachmentFilePath).Output()
				// if err != nil {
				// 	http.Error(w, "Error getting video duration: "+err.Error(), http.StatusInternalServerError)
				// 	return
				// }
				// durationSeconds, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
				// if err != nil {
				// 	http.Error(w, "Error parsing video duration: "+err.Error(), http.StatusInternalServerError)
				// 	return
				// }
				// attachmentRecord.Duration = formatDuration(durationSeconds)

				// Convert video to HLS format
				hlsOutputPath := fmt.Sprintf("/var/www/academyserverapp/hls/videos/%d", time.Now().Unix())
				err = os.MkdirAll(hlsOutputPath, os.ModePerm)
				if err != nil {
					http.Error(w, "Could not create HLS output directory", http.StatusInternalServerError)
					return
				}

				ffmpegCmd := exec.Command("ffmpeg", "-i", attachmentFilePath, "-codec:", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", filepath.Join(hlsOutputPath, "index.m3u8"))
				ffmpegOutput, err := ffmpegCmd.CombinedOutput()
				if err != nil {
					log.Printf("FFmpeg command failed with error: %s\nOutput: %s", err, string(ffmpegOutput))
					http.Error(w, "Could not process video", http.StatusInternalServerError)
					return
				}

				attachmentRecord.FilePath = fmt.Sprintf("/hls/videos/%d/index.m3u8", time.Now().Unix())
			}

			attachments = append(attachments, attachmentRecord)
		}

		// // handle tags
		// tagIDs := r.Form["tag_ids"]
		// for _, tagID := range tagIDs {
		//     var tag tag.Tag
		//     if err := db.First(&tag, tagID).Error; err != nil {
		//         http.Error(w, "Tag not found", http.StatusNotFound)
		//         return
		//     }
		//     content.Tags = append(content.Tags, tag)
		// }

		// Assume tags are sent as an array of strings in the request body
		// tags := make([]string, 0)
		// if err := json.Unmarshal([]byte(r.FormValue("tags")), &tags); err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		// Assume tags are sent as an array of strings in the request body
		var tags []string
		if err := json.Unmarshal([]byte(r.FormValue("tags")), &tags); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//handle categories
		categoryIDs := r.Form["category_ids"] // Assuming category IDs are sent as form values

		// Fetch categories by IDs
		var categories []category.Category
		for _, idStr := range categoryIDs {
			id, _ := strconv.Atoi(idStr)
			var category category.Category
			if err := db.First(&category, id).Error; err != nil {
				http.Error(w, "Category not found", http.StatusNotFound)
				return
			}
			categories = append(categories, category)
		}

		// Store content record in database
		contentRecord := Content{
			Title:        title,
			Description:  description,
			URL:          mainFileURL,
			AuthorID:     uint(AuthorId),
			AuthorName:   AuthorName,
			MainFilePath: mainFileURL,
			MainFileType: filepath.Ext(mainFileHeader.Filename),
			ImageURL:     imageURL,
			Duration:     duration,
			CategoryID:   uint(categoryID),
			CreatedAt:    time.Now(),
			Attachments:  attachments,
			Tags:         tags,
			Categories:   categories,
		}
		if err := db.Create(&contentRecord).Error; err != nil {
			http.Error(w, "Could not save content record", http.StatusInternalServerError)
			return
		}

		// Respond with the content details
		response := UploadResponse{
			ID:          contentRecord.ID,
			ContentURL:  mainFileURL,
			ImageURL:    imageURL,
			Duration:    duration,
			Attachments: attachments,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
