package models

import (
	"gorm.io/gorm"
)


type Excercise struct{
	gorm.Model
	ID					string `json:"id" gorm:"primaryKey"`
	Title				string `json:"title"`
	VideoURL			string `json:"videoUrl"`
	Description 		string `json:"description"`
	Popularity			int	   `json:"popoularity"`
	SubMuscleGroupID	string `json:"subMuscleGroupID"`
}

type SubMuscleGroup struct{
	gorm.Model
	ID				string `json:"id" gorm:"primaryKey"`
	Title			string `json:"title"`
	Excercises		[]Excercise `gorm:"foreignKey:SubMuscleGroupID"`
	MuscleGroupID 	string       `json:"muscleGroupID"`
}

type MuscleGroup struct{
	gorm.Model
	ID					string `json:"id" gorm:"primaryKey"`
	Title				string `json:"title"`
	SubMuscleGroups 	[]SubMuscleGroup `gorm:"foreignKey:MuscleGroupID"`
}