// internal/recommendation/model.go
package recommendation

// RecommendationModel represents the structure of the recommendation in the database
type RecommendationModel struct {
    ID          uint   `gorm:"primary_key"`
    Title       string `gorm:"type:varchar(100);not null"`
    Description string `gorm:"type:text"`
}

// TableName specifies the table name for RecommendationModel
func (RecommendationModel) TableName() string {
    return "recommendations"
}
