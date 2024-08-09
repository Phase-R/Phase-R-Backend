package controllers

import (
	"net/http"

	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/Phase-R/Phase-R-Backend/services/mental_health/db"
	"github.com/gin-gonic/gin"
)

func FetchQuestionSet(ctx *gin.Context) {
	userID, err := ctx.Cookie("user_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in cookies"})
		return
	}

	var qSet models.QuestionSet

	query := `
		SELECT qs.*
		FROM question_sets qs
		LEFT JOIN marked_question_sets mqs ON qs.question_set_id = mqs.question_set_id AND mqs.user_id = ?
		WHERE mqs.user_id IS NULL
		ORDER BY qs.created_at ASC
		LIMIT 1
	`

	// There exists atleast one question set that isn't marked
	res := db.DB.Raw(query, userID).Scan(&qSet)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	// If all questions are marked. Make all of them unmarked and start over
	if qSet.QuestionSetID == "" {
		// unmarking all sets
		res = db.DB.Exec("DELETE FROM marked_question_sets WHERE user_id = ?", userID)
		if res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
			return
		}

		// Fetching the first question set again
		res = db.DB.Raw(query, userID).Scan(&qSet)
		if res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"question_set": qSet})
}

// essentially triggered when answers are submitted 
// func MarkQuestionSet(ctx *gin.Context) {}