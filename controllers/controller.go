package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"upload-files/models"
)

// FileController is a struct that contains a pointer to the database
type FileController struct {
	DB *gorm.DB
}

// UploadFile is a function that handles the upload of a single file
func (c *FileController) UploadFile(ctx *gin.Context) {
	/*
		UploadFile function handles the upload of a single file.
		It gets the file from the form data, saves it to the defined path, 
		generates a unique identifier for the file, saves the file metadata to the database,
		and returns a success message and the file metadata.
	*/
	// Get the file from the form data
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Define the path where the file will be saved
	filePath := filepath.Join("uploads", file.Filename)
	// Save the file to the defined path
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	// Generate a unique identifier for the file
	uuid := uuid.New().String()
	// Save file metadata to database
	fileMetadata := models.File{
		Filename: file.Filename,
		UUID:     uuid,
	}
	if err := c.DB.Create(&fileMetadata).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}
	// Return a success message and the file metadata
	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "Details": fileMetadata})
}

// UploadFiles is a function that handles the upload of multiple files
func (c *FileController) UploadFiles(ctx *gin.Context) {
	/*
		UploadFiles function handles the upload of multiple files.
		It gets the files from the form data, saves each file to the defined path, 
		generates a unique identifier for each file, saves the file metadata to the database,
		and returns a success message and the file metadata.
	*/
	// Get the files from the form data
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	files := form.File["files"]
	var fileModels []models.File
	// Save each file to the defined path and generate a unique identifier for each file
	for _, file := range files {
		filePath := filepath.Join("uploads", file.Filename)
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		fileModels = append(fileModels, models.File{
			UUID:     uuid.New().String(),
			Filename: file.Filename,
		})
	}
	// Save file metadata to database
	err = c.DB.Create(&fileModels).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file information"})
		return
	}
	// Return a success message and the file metadata
	ctx.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"files":   fileModels,
	})
}

// GetFile is a function that retrieves a file from the server
func (c *FileController) GetFile(ctx *gin.Context) {
	/*
		GetFile function retrieves a file from the server.
		It gets the unique identifier of the file to be retrieved, 
		retrieves the file metadata from the database, 
		defines the path of the file to be retrieved, 
		opens the file, reads the first 512 bytes of the file to determine its content type,
		gets the file info, sets the headers for the file transfer, and returns the file.
	*/
	// Get the unique identifier of the file to be retrieved
	uuid := ctx.Param("uuid")
	var file models.File
	// Retrieve the file metadata from the database
	err := c.DB.Where("uuid = ?", uuid).First(&file).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	// Define the path of the file to be retrieved
	filePath := filepath.Join("uploads", file.Filename)
	// Open the file
	fileData, err := os.Open(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileData.Close()
	// Read the first 512 bytes of the file to determine its content type
	fileHeader := make([]byte, 512)
	_, err = fileData.Read(fileHeader)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	fileContentType := http.DetectContentType(fileHeader)
	// Get the file info
	fileInfo, err := fileData.Stat()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}
	// Set the headers for the file transfer and return the file
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	ctx.Header("Content-Type", fileContentType)
	ctx.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	ctx.File(filePath)
}

// DeleteFile is a function that deletes a file from the server and its metadata from the database
func (c *FileController) DeleteFile(ctx *gin.Context) {
	/*
		DeleteFile function deletes a file from the server and its metadata from the database.
		It gets the unique identifier of the file to be deleted, 
		retrieves the file metadata from the database, 
		defines the path of the file to be deleted, 
		deletes the file from the server, 
		deletes the file metadata from the database, 
		and returns a success message.
	*/
	// Get the unique identifier of the file to be deleted
	uuid := ctx.Param("uuid")
	var file models.File
	// Retrieve the file metadata from the database
	err := c.DB.Where("uuid = ?", uuid).First(&file).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	// Define the path of the file to be deleted
	filePath := filepath.Join("uploads", file.Filename)
	// Delete the file from the server
	err = os.Remove(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from upload folder"})
		return
	}
	// Delete the file metadata from the database
	err = c.DB.Delete(&file).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from database"})
		return
	}
	// Return a success message
	ctx.JSON(http.StatusOK, gin.H{
		"message": "File " + file.Filename + " deleted successfully",
	})
}