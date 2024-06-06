package user

import (
	"encoding/json"
	"net/http"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to fetch user data
	// Example:
	user := User{Name: "John Doe", Email: "john@example.com"}
	json.NewEncoder(w).Encode(user)
}
