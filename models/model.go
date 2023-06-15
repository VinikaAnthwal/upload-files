package models

import (
	"gorm.io/gorm"
)

// File is a struct that represents a file in the database.
type File struct {
	gorm.Model // GORM model that contains the ID, CreatedAt, UpdatedAt, and DeletedAt fields
	Filename string `gorm:"not null"` // Filename of the file. Cannot be null.
	UUID     string `gorm:"unique;not null"` // UUID of the file. Must be unique and cannot be null.
}