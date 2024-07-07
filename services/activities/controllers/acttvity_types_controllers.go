package controllers

import (
	"fmt"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/gin-gonic/gin"
	"github.com/nrednav/cuid2"
	"gorm.io/gorm"
	"net/http"
)

func CreateActType(ctx *gin.Context, db *gorm.DB) {
	var actType models.ActivityType
	if err := ctx.ShouldBindJSON(&actType); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := cuid2.Generate()
	if id == "" {
		return
	}

	actType.ID = id

	res := db.Create(&actType)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"yohoo": "activity type created", "activity_type": actType})
}

func GetActType(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var actType models.Activities
	res := db.Raw("SELECT * FROM activitytypes WHERE id = ?", id).Scan(&actType)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"activity_types": actType})
}

func DeleteActType(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var actType models.Activities
	res := db.Where("id = ?", id).First(&actType)

	if res.Error != nil {
		errorMsg := fmt.Sprintf("Activity type with id %s not found", id)
		ctx.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
		return
	}

	del := db.Delete(&actType)
	if del.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": del.Error.Error()})
		return
	}

	successMsg := fmt.Sprintf("Activity type with id %s deleted", actType.Title)
	ctx.JSON(http.StatusOK, gin.H{"yohoo": successMsg})

}

func UpdateActType(ctx *gin.Context, db *gorm.DB) {
	id := ctx.Param("id")
	var actType models.Activities
	var chg models.Activities

	err := ctx.ShouldBindJSON(&chg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := db.Where("id = ?", id).First(&actType)

	if res.Error != nil {
		notfoundmessage := fmt.Sprintf("Activity type with id %s not found", id)
		ctx.JSON(http.StatusNotFound, gin.H{"error": notfoundmessage})
		return
	}

	update := db.Model(&actType).Where("id = ?", id).Updates(chg)
	if update.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": update.Error.Error()})
		return
	}

	successMsg := fmt.Sprintf("Activity type with id %s updated", id)
	ctx.JSON(http.StatusOK, gin.H{"yohoo": successMsg})
}
