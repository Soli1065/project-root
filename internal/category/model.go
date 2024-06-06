// internal/category/model.go
package category

// CategoryModel represents the structure of the category in the database
type CategoryModel struct {
    ID          uint   `gorm:"primary_key"`
    Name        string `gorm:"type:varchar(100);not null"`
    Description string `gorm:"type:text"`
}

// TableName specifies the table name for CategoryModel
func (CategoryModel) TableName() string {
    return "categories"
}
