package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       string `json:"id" gorm:"primaryKey;column:id"`
	Username string `json:"username" gorm:"column:username;unique;not null"`
	Fname    string `json:"fname" gorm:"column:fname;not null"`
	Lname    string `json:"lname" gorm:"column:lname;not null"`
	Email    string `json:"email" gorm:"column:email;unique;not null"`
	Password string `json:"-" gorm:"column:password;not null"`
	Age      int    `json:"age" gorm:"column:age;not null"`
	Access   string `json:"access" gorm:"column:access;not null"`
}
