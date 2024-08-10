package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/Phase-R/Phase-R-Backend/services/mental_health/db"
	"github.com/gin-gonic/gin"
)

func ParseToken(tokenString string) (*jwt.Token, error) {
	secretKey := os.Getenv("SECRET_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func FetchQuestionSet(ctx *gin.Context) {
	// Retrieve token from cookies
	token, err := ctx.Cookie("Auth")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in cookies"})
		return
	}

	// Parse the JWT token
	parsedToken, err := ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Extract the user ID from the token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	fmt.Println(claims)
	if !ok || !parsedToken.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userEmail, ok := claims["iss"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	// Fetch the user ID from the database
	var user models.User
	res := db.DB.Where("email = ?", userEmail).First(&user)
	if res.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	userID := user.ID

	// Query to fetch the first unmarked question set
	var qSet models.QuestionSet
	query := `
		SELECT qs.*
		FROM question_sets qs
		LEFT JOIN marked_question_sets mqs ON qs.id = mqs.id AND mqs.user_id = ?
		WHERE mqs.user_id IS NULL
		ORDER BY qs.created_at ASC
		LIMIT 1
	`

	res = db.DB.Raw(query, userID).Scan(&qSet)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	// If no unmarked question sets found, reset all question sets to unmarked
	if qSet.QuestionSetID == "" {
		res = db.DB.Exec("DELETE FROM marked_question_sets WHERE user_id = ?", userID)
		if res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
			return
		}

		// Fetch the first question set again after resetting
		res = db.DB.Raw(query, userID).Scan(&qSet)
		if res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
			return
		}
	}

	// Return the fetched question set
	ctx.JSON(http.StatusOK, gin.H{"question_set": qSet.Questions, "question_set_id": qSet.QuestionSetID})
}

// func ScoreEvaluation(ctx *gin.Context) {
// 	token, err := ctx.Cookie("Auth")
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in cookies"})
// 		return
// 	}

// 	// Parse the JWT token
// 	parsedToken, err := ParseToken(token)
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		return
// 	}

// 	// Extract the user ID from the token
// 	claims, ok := parsedToken.Claims.(jwt.MapClaims)
// 	fmt.Println(claims)
// 	if !ok || !parsedToken.Valid {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
// 		return
// 	}

// 	userEmail, ok := claims["iss"].(string)
// 	if !ok {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
// 		return
// 	}

// 	// Fetch the user ID from the database
// 	var user models.User
// 	res := db.DB.Where("email = ?", userEmail).First(&user)
// 	if res.Error != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}
// 	userID := user.ID

// 	var answerSet struct {
// 		QuestionSetID string `json:"question_set_id"`
// 		Answers       []bool `json:"answers"`
// 	}
// }