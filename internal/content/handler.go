// internal/content/handler.go
package content

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
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

        if err := DeleteContent(db, uint(id)); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusNoContent)
    }
}
