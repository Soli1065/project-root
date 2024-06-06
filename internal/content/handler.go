package content

import (
	"encoding/json"
	"net/http"
)

func GetContentHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to fetch content data
	// Example:
	content := Content{Title: "Sample Video", Description: "A great video for kids"}
	json.NewEncoder(w).Encode(content)
}
