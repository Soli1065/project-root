// internal/category/handler.go
package category

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
)

// GetAllCategoriesHandler handles the request to retrieve all categories
func GetAllCategoriesHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        categories, err := GetAllCategories(db)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(categories)
    }
}

// GetCategoryByIDHandler handles the request to retrieve a category by its ID
func GetCategoryByIDHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        category, err := GetCategoryByID(db, uint(id))
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        json.NewEncoder(w).Encode(category)
    }
}

// CreateCategoryHandler handles the request to create a new category
func CreateCategoryHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var category Category
        if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if err := CreateCategory(db, &category); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(category)
    }
}

// UpdateCategoryHandler handles the request to update an existing category
func UpdateCategoryHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        var category Category
        if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        category.ID = uint(id)
        if err := UpdateCategory(db, &category); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(category)
    }
}

// DeleteCategoryHandler handles the request to delete a category
func DeleteCategoryHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if err := DeleteCategory(db, uint(id)); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusNoContent)
    }
}
