package config

import (
	"fmt"
	"log"

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

	log.Println("Database migration successful ")
}
