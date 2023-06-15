// Package main is the entry point of the application.
package main

import (
	"log"
	"net/http"

	"upload-files/controllers"
	"upload-files/database"
	"upload-files/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// main is the entry point of the application.
func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database
	err = database.InitDatabase()
	if err != nil {
		log.Fatal("Error initializing database")
	}
	
	// Get the global database instance
	db := database.GlobalDB

	// Log a message indicating that the database has been connected
   log.Println("Connected to database!")

	// Create a new instance of the file controller with the global database instance
	fileController := &controllers.FileController{DB: db}

	// Create a new gin router with default middleware
	r := gin.Default()

	// Define the routes for the file controller
	r.POST("/file", fileController.UploadFile)
	r.POST("/files", fileController.UploadFiles)
	r.GET("/file/:uuid", fileController.GetFile)
	r.DELETE("/file/:uuid", fileController.DeleteFile)

	// Auto-migrate the file model to the database
	err = db.AutoMigrate(&models.File{})
	if err != nil {
		log.Fatal("Error migrating database")
	}

	// Start the server on port 8080
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Error starting server")
	}
}