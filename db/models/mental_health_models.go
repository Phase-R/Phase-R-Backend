package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type QuestionSet struct {
	gorm.Model
	ID					int32			`json:"id" gorm:"primaryKey;autoIncrement:true;unique"`
	Questions			pq.StringArray	`json:"questions" gorm:"type:text[]"`
}

type MarkedQuestionSet struct {
	gorm.Model
	UserID				string		`json:"uid" gorm:"uniqueIndex"`
	Marked				pq.Int32Array		`json:"marked" gorm:"type:int[]"`
}