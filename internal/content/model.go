// internal/content/model.go
package content

import (
	"time"
)

// ContentModel represents the structure of the content in the database
type ContentModel struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	CategoryID  uint      `json:"category_id"`
	AuthorID    uint      `json:"author_id"`
	AuthorName  string    `json:"author_name"`
	ImageURL    string    `json:"image_url"`
	ViewCount   int       `json:"view_count"`
	Duration    float64   `json:"duration"`
	IsLive      bool      `json:"is_live"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName specifies the table name for ContentModel
func (ContentModel) TableName() string {
	return "contents"
}
