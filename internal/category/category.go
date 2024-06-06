// internal/category/category.go
package category

import (
	"gorm.io/gorm"
)

// Category represents a category item
type Category struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetAllCategories retrieves all categories from the database
func GetAllCategories(db *gorm.DB) ([]Category, error) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByID retrieves a category by its ID from the database
func GetCategoryByID(db *gorm.DB, id uint) (Category, error) {
	var category Category
	if err := db.First(&category, id).Error; err != nil {
		return category, err
	}
	return category, nil
}

// CreateCategory creates a new category in the database
func CreateCategory(db *gorm.DB, category *Category) error {
	if err := db.Create(category).Error; err != nil {
		return err
	}
	return nil
}

// UpdateCategory updates an existing category in the database
func UpdateCategory(db *gorm.DB, category *Category) error {
	if err := db.Save(category).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCategory deletes a category from the database
func DeleteCategory(db *gorm.DB, id uint) error {
	if err := db.Delete(&Category{}, id).Error; err != nil {
		return err
	}
	return nil
}
