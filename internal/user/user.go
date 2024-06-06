// internal/user/user.go
package user

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(router *mux.Router, db *gorm.DB) {
	router.HandleFunc("/api/users", GetAllUsersHandler(db)).Methods("GET")
	router.HandleFunc("/api/users/{id}", GetUserHandler(db)).Methods("GET")
	router.HandleFunc("/api/users", CreateUserHandler(db)).Methods("POST")
	router.HandleFunc("/api/users/{id}", UpdateUserHandler(db)).Methods("PUT")
	router.HandleFunc("/api/users/{id}", DeleteUserHandler(db)).Methods("DELETE")
}
