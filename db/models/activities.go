package models

import (
	"gorm.io/gorm"
)

// Drill struct represents a drill under a subtype of an activity.
// Embedded GORM struct. Includes ID,CreatedAt,UpdatedAt,DeletedAt fields.

type Drill struct {
	gorm.Model
	ID             string `json:"id" gorm:"primaryKey"`
	ActivityTypeID string `json:"activityTypeID"`
	Title          string `json:"title"`
	Details        string `json:"description"` // Reps, extra info etc.
	VideoURL       string `json:"videoUrl"`
}

type ActivityType struct {
	gorm.Model
	ID          string  `json:"id" gorm:"primaryKey"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Headline    string  `json:"headline"` // Headlines like GET READY FOR etc.
	ImageURL    string  `json:"imageUrl"` // Cover Image
	ActivityID  string  `json:"activityID"`
	Drills      []Drill `gorm:"foreignKey:ActivityTypeID"`
}

type Activities struct {
	gorm.Model
	ID          string         `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Types       []ActivityType `gorm:"foreignKey:ActivityID"`
	ImageURL    string         `json:"imageUrl"`   // Cover Image
	ColourCode  string         `json:"colourCode"` // Sport Colour code
}

func (Activities) TableName() string {
	return "activities"
}
