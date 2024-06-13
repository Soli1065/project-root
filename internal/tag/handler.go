package tag

// import (
// 	"encoding/json"
// 	"net/http"

// 	"strconv"

// 	"github.com/gorilla/mux"
// 	"gorm.io/gorm"
// 	// "os"
// )

// // AddTagsToContentHandler adds tags to a content.
// func AddTagsToContentHandler(db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		contentID, err := strconv.ParseUint(mux.Vars(r)["contentID"], 10, 64)
// 		if err != nil {
// 			http.Error(w, "Invalid content ID", http.StatusBadRequest)
// 			return
// 		}

// 		var tags []Tag
// 		if err := json.NewDecoder(r.Body).Decode(&tags); err != nil {
// 			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
// 			return
// 		}

// 		// Associate each tag with the content ID
// 		for i := range tags {
// 			tags[i].ContentID = uint(contentID)
// 		}

// 		// Save tags to the database
// 		for _, tag := range tags {
// 			if err := CreateTag(db, &tag); err != nil {
// 				http.Error(w, "Failed to create tag", http.StatusInternalServerError)
// 				return
// 			}
// 		}

// 		// Respond with success message or tags
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(tags)
// 	}
// }

// // GetTagsForContentHandler fetches all tags associated with a content.
// func GetTagsForContentHandler(db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		contentID, err := strconv.ParseUint(mux.Vars(r)["contentID"], 10, 64)
// 		if err != nil {
// 			http.Error(w, "Invalid content ID", http.StatusBadRequest)
// 			return
// 		}

// 		tags, err := GetTagsForContent(db, uint(contentID))
// 		if err != nil {
// 			http.Error(w, "Failed to get tags for content", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(tags)
// 	}
// }

// func CreateTag(db *gorm.DB, tag *Tag) error {
// 	return db.Create(tag).Error
// }

// func GetTagsForContent(db *gorm.DB, contentID uint) ([]Tag, error) {
// 	var tags []Tag
// 	err := db.Where("content_id = ?", contentID).Find(&tags).Error
// 	return tags, err
// }

// func UpdateTag(db *gorm.DB, tag *Tag) error {
// 	return db.Save(tag).Error
// }

// func DeleteTag(db *gorm.DB, tagID uint) error {
// 	return db.Delete(&Tag{}, tagID).Error
// }
