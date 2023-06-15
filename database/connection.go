package database

import (
	"fmt"
	"log"
	"upload-files/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GlobalDB is a global variable that holds the database connection
var GlobalDB *gorm.DB

// InitDatabase initializes the database connection and performs auto-migration
func InitDatabase() (err error) {
	// Read the environment variables from .env file
	config, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error reading .env file")
	}

	// Create a Data Source Name (DSN) string to connect to the database
	dsn := fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config["DB_USERNAME"],
		config["DB_PASSWORD"],
		config["DB_HOST"],
		config["DB_DATABASE"],
	)

	// Open a connection to the database using the DSN string
	GlobalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	// Perform auto-migration to create the necessary tables in the database
	err = GlobalDB.AutoMigrate(&models.File{})
	if err != nil {
		return
	}

	return
}