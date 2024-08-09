package controllers

import (
	"net/http"

	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/Phase-R/Phase-R-Backend/services/mental_health/db"
	"github.com/gin-gonic/gin"
)

func FetchQuestionSet(ctx *gin.Context) {
	id := ctx.Param("set_no")
	var qSet models.QuestionSet

	res := db.DB.Raw("SELECT * FROM question_sets WHERE id = ?", id).Scan(&qSet)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"question_set": qSet})
}