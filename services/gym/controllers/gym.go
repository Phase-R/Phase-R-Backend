package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/Phase-R/Phase-R-Backend/services/gym/db" 
)

func GetAllGymData(ctx *gin.Context) {
	var muscleGroups []models.MuscleGroup

	// Preload related SubMuscleGroups and their associated Excercises
	if err := db.DB.Preload("SubMuscleGroups.Excercises").Find(&muscleGroups).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	ctx.JSON(http.StatusOK, muscleGroups)
}
