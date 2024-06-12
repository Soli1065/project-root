// internal/content/model.go
package content

import (
	"project-root/internal/attachment"
	"time"
)

// ContentModel represents the structure of the content in the database
type ContentModel struct {
	ID           uint                    `gorm:"primaryKey" json:"id"`
	Title        string                  `json:"title"`
	Description  string                  `json:"description"`
	URL          string                  `json:"url"`
	CategoryID   uint                    `json:"category_id"`
	AuthorID     uint                    `json:"author_id"`
	AuthorName   string                  `json:"author_name"`
	ImageURL     string                  `json:"image_url"`
	ViewCount    int                     `json:"view_count"`
	Duration     string                  `json:"duration"`
	IsLive       bool                    `json:"is_live"`
	CreatedAt    time.Time               `json:"created_at"`
	MainFilePath string                  `json:"main_file_path"`
	MainFileType string                  `json:"main_file_type"`
	Attachments  []attachment.Attachment `json:"attachments"`
}

// TableName specifies the table name for ContentModel
func (ContentModel) TableName() string {
	return "contents"
}
