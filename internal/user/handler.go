// internal/user/handler.go
package user

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    "dastanamon/internal/user/model"
    "strconv"
)

func GetAllUsersHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var users []model.User
        if err := db.Find(&users).Error; err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    }
}

func GetUserHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }

        var user model.User
        if err := db.First(&user, id).Error; err != nil {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    }
}

func CreateUserHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user model.User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if err := db.Create(&user).Error; err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
    }
}

func UpdateUserHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }

        var user model.User
        if err := db.First(&user, id).Error; err != nil {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if err := db.Save(&user).Error; err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    }
}

func DeleteUserHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }

        if err := db.Delete(&model.User{}, id).Error; err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}
