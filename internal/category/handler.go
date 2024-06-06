package category

import (
	"encoding/json"
	"net/http"
)

func GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to fetch category data
	// Example:
	category := Category{Name: "Cartoons"}
	json.NewEncoder(w).Encode(category)
}
