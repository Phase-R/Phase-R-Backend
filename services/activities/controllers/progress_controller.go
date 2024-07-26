package controllers

import (
    "net/http"
    "github.com/Phase-R/Phase-R-Backend/activities/db"
    "github.com/Phase-R/Phase-R-Backend/db/models"
    "github.com/gin-gonic/gin"
    "github.com/nrednav/cuid2"
    "gorm.io/gorm"
)


func GetUserProgress(ctx *gin.Context) {
    userID := ctx.Param("userId")
    var userActivityMappings []models.UserActivityMapping

    err := db.DB.Preload("ActivityMappings.DrillMappings").Where("user_id = ?", userID).Find(&userActivityMappings).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"message": "User progress not found"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"data": userActivityMappings})
}



func UpdateUserActivityCompletion(ctx *gin.Context) {
    id := ctx.Param("userActivityId")
    var body struct {
        Completed bool `json:"completed"`
    }

    if err := ctx.ShouldBindJSON(&body); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var userActivityMapping models.UserActivityMapping
    res := db.DB.Where("id = ?", id).First(&userActivityMapping)
    if res.Error != nil {
        if res.Error == gorm.ErrRecordNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"message": "User activity mapping not found"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
        return
    }

    userActivityMapping.Completed = body.Completed
    db.DB.Save(&userActivityMapping)

    ctx.JSON(http.StatusOK, gin.H{"message": "User activity completion updated successfully"})
}


func UpdateDrillCompletion(ctx *gin.Context) {
    id := ctx.Param("drillId")
    var body struct {
        Completed bool `json:"completed"`
    }

    if err := ctx.ShouldBindJSON(&body); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var drillCompletion models.DrillCompletion
    res := db.DB.Where("id = ?", id).First(&drillCompletion)
    if res.Error != nil {
        if res.Error == gorm.ErrRecordNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"message": "Drill not found"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
        return
    }

    drillCompletion.Completed = body.Completed
    db.DB.Save(&drillCompletion)

    ctx.JSON(http.StatusOK, gin.H{"message": "Drill completion updated successfully"})
}
