package models

import (
	"gorm.io/gorm"
)

type QuestionSet struct {
	gorm.Model
	QuestionSetID		string
	Questions			[]string
}

type MarkedQuestionSet struct {
	gorm.Model
	UserID				string					
	Marked				[]bool
}