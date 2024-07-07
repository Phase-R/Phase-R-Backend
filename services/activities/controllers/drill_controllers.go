package controllers

import (
	"fmt"
	"net/http"

	"github.com/Phase-R/Phase-R-Backend/activities/db"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/gin-gonic/gin"
	"github.com/nrednav/cuid2"
)

func CreateDrill(ctx *gin.Context) {
	var drill models.Drill
	if err := ctx.ShouldBindJSON(&drill); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := cuid2.Generate()
	if id == "" {
		return
	}

	drill.ID = id

	res := db.DB.Create(&drill)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Drill created successfully."})
}

func GetDrillsByType(ctx *gin.Context) {
	type_id := ctx.Param("typeid")

	var activityType models.ActivityType
	err := db.DB.Preload("Drills").First(&activityType, "id = ?", type_id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if len(activityType.Drills) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No drills found!"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": activityType.Drills})

}

func DeleteDrill(ctx *gin.Context) {
	id := ctx.Param("id")
	var drill models.Drill
	res := db.DB.Where("id = ?", id).First(&drill)

	if res.Error != nil {
		errorMsg := fmt.Sprintf("No drills with id: %s were found!", id)
		ctx.JSON(http.StatusNotFound, gin.H{"message": errorMsg})
		return
	}

	del := db.DB.Delete(&drill)
	if del.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": del.Error.Error()})
		return
	}

	successMsg := fmt.Sprintf("Deleted drill: %s successfully!", drill.Title)
	ctx.JSON(http.StatusOK, gin.H{"message": successMsg})
}

func GetDrill(ctx *gin.Context) {
	id := ctx.Param("id")
	var drill models.Drill
	res := db.DB.Where("id = ?", id).First(&drill)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": drill})
}

func UpdateDrill(ctx *gin.Context) {
	id := ctx.Param("id")
	var drill models.Drill
	var changes models.Drill

	err := ctx.ShouldBindJSON(&changes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := db.DB.Where("id = ?", id).First(&drill)

	if res.Error != nil {
		notFoundMsg := fmt.Sprintf("No drills with the id: %s were found!", id)
		ctx.JSON(http.StatusNotFound, gin.H{"message": notFoundMsg})
		return
	}

	update := db.DB.Model(&drill).Where("id = ?", id).Updates(changes)
	if update.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": update.Error.Error()})
		return
	}

	successMsg := fmt.Sprintf("Updated drill with id: %s successfully!", id)
	ctx.JSON(http.StatusOK, gin.H{"message": successMsg})
}
