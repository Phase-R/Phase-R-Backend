package models

import (
	"gorm.io/gorm"
)

type Drill struct {
	gorm.Model
	ID          string `json:"id" gorm:"column:id"`
	Type        string `json:"type" gorm:"column:type_act"`
	Title       string `json:"title" gorm:"column:title"`
	Description string `json:"description" gorm:"column:description"`
	VideoURL    string `json:"videoUrl" gorm:"column:videourl"`
}
