package controllers

import (
	"net/http"
	"time"

	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/Phase-R/Phase-R-Backend/services/mental_health/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	// "gorm.io/gorm"
)

func PostThoughts(ctx *gin.Context) {
    token, err := ctx.Cookie("Auth")
    if err != nil {
        ctx.JSON(400, gin.H{"message": "Auth cookie not found"})
        return
    }

    // Parse the JWT token
    parsedToken, err := ParseToken(token)
    if err != nil {
        ctx.JSON(401, gin.H{"message": "Invalid token"})
        return
    }

    // Extract the user ID from the token
    claims, ok := parsedToken.Claims.(jwt.MapClaims)
    if !ok || !parsedToken.Valid {
        ctx.JSON(403, gin.H{"message": "Invalid token claims"})
        return
    }

    userEmail, ok := claims["iss"].(string) // Assuming "iss" is used for email
    if !ok {
        ctx.JSON(403, gin.H{"message": "Invalid user ID in token"})
        return
    }

    // Fetch the user ID from the database
    var user models.User
    res := db.DB.Where("email = ?", userEmail).First(&user)
    if res.Error != nil {
        ctx.JSON(404, gin.H{"message": "User not found"})
        return
    }
    userID := user.ID

    var body struct {
        Thought string `json:"thought"`
    }

    if err := ctx.BindJSON(&body); err != nil {
        ctx.JSON(400, gin.H{"error": "Failed to read the body"})
        return
    }

    // Get the current date (without time)
    currentDate := time.Now().Truncate(24 * time.Hour)

    var existingThought models.Thoughts
    result := db.DB.Where("user_id = ? AND date_of_thought >= ? AND date_of_thought < ?", 
        userID, currentDate, currentDate.Add(24*time.Hour)).First(&existingThought)

    if result.Error == nil {
        // Thought for today exists, update it
        existingThought.Thought = body.Thought
        if err := db.DB.Save(&existingThought).Error; err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        ctx.JSON(http.StatusOK, gin.H{"message": "Thought updated successfully"})
    } else {
        // Create a new thought
        newThought := models.Thoughts{
            DateOfThought: currentDate,
            UserID:        userID,
            Thought:       body.Thought,
        }
        if err := db.DB.Create(&newThought).Error; err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        ctx.JSON(http.StatusCreated, gin.H{"message": "New thought added successfully"})
    }
}
