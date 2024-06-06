// internal/recommendation/recommendation.go
package recommendation

import (
	"gorm.io/gorm"
)

// Recommendation represents a recommendation item
type Recommendation struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// GetAllRecommendations retrieves all recommendations from the database
func GetAllRecommendations(db *gorm.DB) ([]Recommendation, error) {
	var recommendations []Recommendation
	if err := db.Find(&recommendations).Error; err != nil {
		return nil, err
	}
	return recommendations, nil
}

// GetRecommendationByID retrieves a recommendation by its ID from the database
func GetRecommendationByID(db *gorm.DB, id uint) (Recommendation, error) {
	var recommendation Recommendation
	if err := db.First(&recommendation, id).Error; err != nil {
		return recommendation, err
	}
	return recommendation, nil
}

// CreateRecommendation creates a new recommendation in the database
func CreateRecommendation(db *gorm.DB, recommendation *Recommendation) error {
	if err := db.Create(recommendation).Error; err != nil {
		return err
	}
	return nil
}

// UpdateRecommendation updates an existing recommendation in the database
func UpdateRecommendation(db *gorm.DB, recommendation *Recommendation) error {
	if err := db.Save(recommendation).Error; err != nil {
		return err
	}
	return nil
}

// DeleteRecommendation deletes a recommendation from the database
func DeleteRecommendation(db *gorm.DB, id uint) error {
	if err := db.Delete(&Recommendation{}, id).Error; err != nil {
		return err
	}
	return nil
}
