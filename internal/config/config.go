package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize initializes the database connection
func Initialize() {
	// Database connection details
	dbHost := "localhost"
	dbPort := "5432"
	dbUser := "postgres"
	dbPassword := "Soli@123456789"
	dbName := "academydb"

	// Create the Data Source Name (DSN)
	dsn := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"

	fmt.Printf("This is dsn:  %v ", dsn)

	// Initialize the database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
}
