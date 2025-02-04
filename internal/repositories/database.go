package repositories

import (
	"fmt"
	"go_swift/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a global variable to hold the database connection
var DB *gorm.DB

// InitDB initializes the database connection and performs migrations
func InitDB() {
	// Print environment variables for debugging purposes
	fmt.Println("🚀 Debugging Environment Variables:")
	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))
	fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))

	// Create the Data Source Name (DSN) for PostgreSQL connection
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Open a connection to the PostgreSQL database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}
	DB = db

	// Perform database migrations for the SwiftCode model
	err = DB.AutoMigrate(&models.SwiftCode{})
	if err != nil {
		log.Fatalf("❌ Database migration failed: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL and migrated database successfully!")
}
