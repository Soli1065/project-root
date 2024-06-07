// internal/video/model.go
package video

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string `gorm:"type:text"`
	URL         string `gorm:"not null;unique"`
	IsLive      bool   `gorm:"default:false"`
}
