package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type QuestionSet struct {
	gorm.Model
	QuestionSetID		string 			`json:"id" gorm:"primaryKey"`
	Questions			pq.StringArray	`json:"questions" gorm:"type:text"`
}

type MarkedQuestionSet struct {
	gorm.Model
	ID					string		`json:"id" gorm:"primaryKey"`
	UserID				string		`json:"uid" gorm:"foreignKey:"`
	Marked				[]bool		`json:"marked" gorm:"type:boolean[]"`
}