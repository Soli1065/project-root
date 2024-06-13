package tag

import (
	"project-root/internal/content"
	"time"
)

type Tag struct {
	ID        uint              `gorm:"primary_key" json:"id"`
	Name      string            `json:"name"`
	CreatedAt time.Time         `gorm:"autoCreateTime"`
	Contents  []content.Content `json:"contents" gorm:"many2many:content_tags;"`
}
