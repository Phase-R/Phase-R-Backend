package models

import (
    "fmt"
    "github.com/Phase-R/Phase-R-Backend/auth/utils"
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
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
    if u.OTP != "" {
        hashedOTP, err := utils.PwdSaltAndHash(u.OTP)
        if err != nil {
            return fmt.Errorf("failed to hash OTP: %v", err)
        }
        u.OTP = hashedOTP
        if err := tx.Save(u).Error; err != nil {
            return fmt.Errorf("failed to save hashed OTP: %v", err)
        }
    }
    return nil
}
