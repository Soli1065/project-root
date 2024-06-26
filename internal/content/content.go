package content

import (
	"project-root/internal/attachment"
	"project-root/internal/category"
	"project-root/internal/comment"

	// "project-root/internal/tag"

	"time"

	"gorm.io/gorm"

	"github.com/lib/pq" // Import pq package for PostgreSQL array support
)

// Content represents a content record in the database
type Content struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	// CategoryID   uint                    `json:"category_id"`
	AuthorID          uint                    `json:"author_id"`
	AuthorName        string                  `json:"author_name"`
	CreatedAt         time.Time               `gorm:"autoCreateTime" json:"created_at"`
	ImageURL          string                  `gorm:"type:varchar(255)" json:"image_url"`
	ViewCount         uint                    `gorm:"default:0" json:"view_count"`
	LikesCount        uint                    `gorm:"default:0" json:"likes_count"` // New field for likes count
	Duration          string                  `json:"duration"`
	IsLive            bool                    `gorm:"default:false" json:"is_live"`
	MainFilePath      string                  `json:"main_file_path"`
	MainFileType      string                  `json:"main_file_type"`
	CategoryIDs       []uint                  `gorm:"-" json:"category_ids"`                           // Ignored by GORM, used for handling IDs from request
	Categories        []category.Category     `gorm:"many2many:content_categories;" json:"categories"` // Define the many-to-many relationship
	Attachments       []attachment.Attachment `json:"attachments" gorm:"foreignKey:ContentID"`
	Comments          []comment.Comment       `json:"comments" gorm:"foreignKey:ContentID"`
	Tags              pq.StringArray          `gorm:"type:text[]"`                               // PostgreSQL array type for tags
	IsApproved        bool                    `gorm:"default:false" json:"is_approved"`          // New field for approval status
	RelatedContentIDs pq.Int64Array           `gorm:"type:integer[]" json:"related_content_ids"` // Use pq.Int64Array for PostgreSQL integer arrays

}

// GetAllContents retrieves all contents from the database
func GetAllContents(db *gorm.DB) ([]Content, error) {
	var contents []Content
	if err := db.Find(&contents).Error; err != nil {
		return nil, err
	}
	return contents, nil
}

// GetContentByID retrieves a content by its ID from the database
func GetContentByID(db *gorm.DB, id uint) (Content, error) {
	var content Content
	if err := db.First(&content, id).Error; err != nil {
		return content, err
	}
	return content, nil
}

// CreateContent creates a new content item in the database
func CreateContent(db *gorm.DB, content *Content) error {
	if err := db.Create(content).Error; err != nil {
		return err
	}
	return nil
}

// UpdateContent updates an existing content item in the database
func UpdateContent(db *gorm.DB, content *Content) error {
	if err := db.Save(content).Error; err != nil {
		return err
	}
	return nil
}

// DeleteContent deletes an existing content item in the database
func DeleteContent(db *gorm.DB, content *Content) error {
	if err := db.Delete(content).Error; err != nil {
		return err
	}
	return nil
}
