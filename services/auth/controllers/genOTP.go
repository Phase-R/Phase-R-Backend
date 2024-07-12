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

func GenerateOTP(c *gin.Context) (string, error) {
    var body struct {
        Email string `json:"email"`
    }

    if c.Bind(&body) != nil {
        return "", fmt.Errorf("invalid request body")
    }

    var user models.User
    result := db.DB.Where("email = ?", body.Email).First(&user)
    if result.Error != nil {
        return "", fmt.Errorf("email not found")
    }

    // Generate a 6-digit OTP
    otp := rand.Intn(1000000)

    // Hash the OTP using the salted hash function
    hashedOTP, err := utils.PwdSaltAndHash(fmt.Sprintf("%d", otp))
    if err != nil {
        return "", fmt.Errorf("failed to hash OTP")
    }

    // Store the hashed OTP
    user.OTP = hashedOTP

    // Save the user with the new OTP
    db.DB.Save(&user)

    // Return the plain OTP for sending via email
    return fmt.Sprintf("%06d", otp), nil
}
