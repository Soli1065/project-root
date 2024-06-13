package config

import (
	"fmt"
	"log"

	"project-root/internal/attachment"
	"project-root/internal/category"
	"project-root/internal/content"
	"project-root/internal/recommendation"
	"project-root/internal/user"
	"project-root/internal/video"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize initializes the database connection
func Initialize() {
	// Database connection details
	// dbHost := "localhost"
	// dbPort := "5432"
	// dbUser := "postgres"
	// dbPassword := "Soli@123456789"
	// dbName := "academydb"

	dbHost := "localhost"
	dbPort := "5432"
	dbUser := "academydbuser"
	dbPassword := "bpj12345"
	dbName := "academydb"

	// Create the Data Source Name (DSN)
	dsn := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"

	fmt.Printf("This is dsn:  %v \n", dsn)

	// Initialize the database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Perform migrations for each service

	// Migrate the schema
	// DB.AutoMigrate(&user.User{}, &category.CategoryModel{}, &content.Content{}, &recommendation.Recommendation{}, &video.Video{})

	if err := DB.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("Failed to migrate user model: %v", err)
	}

	if err := DB.AutoMigrate(&content.Content{}); err != nil {
		log.Fatalf("Failed to migrate content model: %v", err)
	}

	if err := DB.AutoMigrate(&category.CategoryModel{}); err != nil {
		log.Fatalf("Failed to migrate category model: %v", err)
	}

	if err := DB.AutoMigrate(&recommendation.Recommendation{}); err != nil {
		log.Fatalf("Failed to migrate recommendation model: %v", err)
	}

	if err := DB.AutoMigrate(&video.Video{}); err != nil {
		log.Fatalf("Failed to migrate video model: %v", err)
	}

	if err := DB.AutoMigrate(&attachment.Attachment{}); err != nil {
		log.Fatalf("Failed to migrate attachment model: %v", err)
	}

	// if err := DB.AutoMigrate(&comment.Comment{}); err != nil {
	// 	log.Fatalf("Failed to migrate attachment model: %v", err)
	// }

	log.Println("Database migration successful ")
}

// Function to insert new content into the database
// func InsertContent(content content.Content) (int, error) {
// 	result := DB.Create(&content)
// 	if result.Error != nil {
// 		return 0, result.Error
// 	}
// 	return int(content.ID), nil
// }

// Function to insert new attachment into the database
// func InsertAttachment(attachment Attachment) (int, error) {
// 	result := DB.Create(&attachment)
// 	if result.Error != nil {
// 		return 0, result.Error
// 	}
// 	return int(attachment.ID), nil
// }
