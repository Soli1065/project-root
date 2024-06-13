package attachment

import "time"

// Attachment represents an attachment record in the database
type Attachment struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	ContentID uint      `json:"content_id"`
	FilePath  string    `json:"file_path"`
	FileType  string    `json:"file_type"`
	CreatedAt time.Time `json:"created_at"`
}
