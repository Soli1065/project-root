package attachment

import "time"

type Attachment struct {
	ID        int       `json:"id"`
	ContentID int       `json:"content_id"`
	FilePath  string    `json:"file_path"`
	FileType  string    `json:"file_type"`
	CreatedAt time.Time `json:"created_at"`
}