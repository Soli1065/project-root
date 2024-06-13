package comment

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	ContentID uint      `json:"content_id"`
	UserID    uint      `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
