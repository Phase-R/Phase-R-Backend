package controllers

// package main

import (
	"github.com/Phase-R/Phase-R-Backend/auth/utils"
	"github.com/Phase-R/Phase-R-Backend/auth/db"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/gin-gonic/gin"
	"github.com/nrednav/cuid2"
	"log"
	"net/http"
)



func CreateUser(ctx *gin.Context) {
	const uniqueViolation = "23505"

	var newUser models.User

	id := cuid2.Generate()
	if id == "" {
		return
	}

	newUser.ID = id

	hash, err := utils.PwdSaltAndHash(newUser.Password)
	if err != nil {
		log.Fatal("could not hash password", err)
	}

	newUser.Password = hash

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := db.DB.Create(&newUser)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"yohoo": "new user created."})
}

func FetchUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User
	res := db.DB.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&user)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var updatedUser models.User

	if err := ctx.ShouldBindJSON(&updatedUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := db.DB.Model(&models.User{}).Where("id = ?", id).Updates(updatedUser)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	res := db.DB.Delete(&models.User{}, id)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
