package tag

// import "time"

// // Tag represents a tag associated with a content.
// type Tag struct {
// 	ID         uint   `gorm:"primary_key" json:"id"`
// 	Name       string `gorm:"unique;not null" json:"name"`
// 	ContentID  uint   `json:"content_id"` // Foreign key to associate with Content
// 	CreatedAt  time.Time
// 	UpdatedAt  time.Time
// }

// Tag represents a tag associated with content.
// type Tag struct {
//     ID        uint   `gorm:"primary_key" json:"id"`
//     Name      string `json:"name"`
//     CreatedAt time.Time `json:"created_at"`
//     UpdatedAt time.Time `json:"updated_at"`
//     DeletedAt *time.Time `sql:"index" json:"-"`
// }

// // ContentTag represents the association between content and tags.
// type ContentTag struct {
//     ContentID uint `gorm:"primaryKey" json:"content_id"`
//     TagID     uint `gorm:"primaryKey" json:"tag_id"`
// }
