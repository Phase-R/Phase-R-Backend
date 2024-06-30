package models

import (
	"gorm.io/gorm"
)

// User struct represents a user in the system.
// Embedded GORM struct. Includes ID,CreatedAt,UpdatedAt,DeletedAt fields.
type User struct {
	gorm.Model
	ID       string `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Fname    string `json:"fname" gorm:"not null"`
	Lname    string `json:"lname" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Age      int    `json:"age" gorm:"not null"`
	Access   string `json:"access" gorm:"not null"` // free or premium
}
