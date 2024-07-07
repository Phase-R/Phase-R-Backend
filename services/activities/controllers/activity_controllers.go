package controllers

import (
	"fmt"
	"net/http"

	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/gin-gonic/gin"
	"github.com/nrednav/cuid2"
	"gorm.io/gorm"
)

func CreateActivity(ctx *gin.Context, db *gorm.DB) {
	var activity models.Activities
	if err := ctx.ShouldBindJSON(&activity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := cuid2.Generate()
	activity.ID = id

	res := db.Create(&activity)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "activity created successfully", "activity": activity})
}

func GetActivity(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var activity models.Activities
	err := db.Preload("Types").First(&activity, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Activity not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": activity})
}

func UpdateActivity(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var changes models.Activities

	if err := ctx.ShouldBindJSON(&changes); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var activity models.Activities
	res := db.Where("id = ?", id).First(&activity)
	if res.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No activity with the ID: %s found", id)})
		return
	}

	update := db.Model(&activity).Updates(changes)
	if update.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": update.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Updated activity with ID: %s successfully", id), "activity": activity})
}

func DeleteActivity(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var activity models.Activities
	res := db.Where("id = ?", id).First(&activity)
	if res.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No activity with ID: %s found", id)})
		return
	}

	del := db.Delete(&activity)
	if del.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": del.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Deleted activity: %s successfully", activity.Title)})
}
