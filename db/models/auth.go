package models

import (
	"gorm.io/gorm"
	"time"
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
	Access   string `json:"access" gorm:"not null"` //free or premium
	Verified bool   `json:"verified" gorm:"not null"`
	OTP      string `json:"otp"`
}


//afterupdate gorm hook to auto delete all OTPs every 5 mins
func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	// Schedule a job to reset the OTP after 5 minutes
	go func() {
		time.Sleep(5 * time.Minute)
		tx.Model(&User{}).Where("id = ?", u.ID).Update("otp", nil)
	}()
	return
}