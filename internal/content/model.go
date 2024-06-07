// internal/content/model.go
package content

import (
	"time"
)

// ContentModel represents the structure of the content in the database
type ContentModel struct {
	ID          uint   `gorm:"primary_key"`
	Title       string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text"`
	URL         string `gorm:"type:varchar(255);not null"`
	CategoryID  uint
	AuthorID    uint      `gorm:"not null"`
	AuthorName  string    `gorm:"type:varchar(100);not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

// TableName specifies the table name for ContentModel
func (ContentModel) TableName() string {
	return "contents"
}
