// internal/content/model.go
package content

// ContentModel represents the structure of the content in the database
type ContentModel struct {
    ID          uint   `gorm:"primary_key"`
    Title       string `gorm:"type:varchar(100);not null"`
    Description string `gorm:"type:text"`
    URL         string `gorm:"type:varchar(255);not null"`
    CategoryID  uint
}

// TableName specifies the table name for ContentModel
func (ContentModel) TableName() string {
    return "contents"
}
