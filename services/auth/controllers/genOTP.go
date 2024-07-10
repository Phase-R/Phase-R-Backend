package controllers

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/Phase-R/Phase-R-Backend/auth/db"
	"github.com/Phase-R/Phase-R-Backend/auth/utils"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/gin-gonic/gin"
)

// GenerateOTP handles the OTP generation and emailing process
func GenerateOTP(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var user models.User
	result := db.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
		return
	}

	// Generate a 6-digit OTP
	otp := rand.Intn(1000000)

	// Hash the OTP using the salted hash function
	hashedOTP, err := utils.PwdSaltAndHash(fmt.Sprintf("%d", otp))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash OTP"})
		return
	}

	// Store the hashed OTP
	user.OTP = hashedOTP
	db.DB.Save(&user)

	// Send plain OTP via email
	sendEmail(user.Email, fmt.Sprintf("%06d", otp))

	c.JSON(http.StatusOK, gin.H{"message": "password reset OTP sent"})
}
