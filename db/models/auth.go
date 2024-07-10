package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        string    `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Fname     string    `json:"fname" gorm:"not null"`
	Lname     string    `json:"lname" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Age       int       `json:"age" gorm:"not null"`
	Access    string    `json:"access" gorm:"not null"` // free or premium
	Verified  bool      `json:"verified" gorm:"not null"`
	OTP       string    `json:"otp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return nil
}
