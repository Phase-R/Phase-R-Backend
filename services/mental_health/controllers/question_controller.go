package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/Phase-R/Phase-R-Backend/services/mental_health/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	// "github.com/lib/pq"
	"gorm.io/gorm"
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

func AddQuestionSet(ctx *gin.Context) {
	token, err := ctx.Cookie("Auth")
	// token, err := ctk.Cooker("Admin")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in cookies"})
		return
	}

	parsedToken, err := ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}
	userEmail, ok := claims["iss"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}
	fmt.Println(userEmail)

	var questionSet models.QuestionSet

	if ctx.Bind(&questionSet) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read the body"})
		return
	}

	res := db.DB.Create(&questionSet)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Question Set created successfully"})
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
		SELECT qs.id, qs.questions
	FROM question_sets qs
	LEFT JOIN marked_question_sets mqs ON mqs.user_id = ?
	WHERE qs.id NOT IN (
		SELECT ALL(unnest(mqs.marked))
		FROM marked_question_sets mqs
		WHERE mqs.user_id = ?
	)
	LIMIT 1
	`

	res = db.DB.Raw(query, userID, userID).Scan(&qSet)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	// Return the fetched question set
	ctx.JSON(http.StatusOK, gin.H{"question_set": qSet.Questions, "question_set_id": qSet.ID})
}

func ScoreEvaluation(ctx *gin.Context) {
	token, err := ctx.Cookie("Auth")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in cookies"})
		return
	}

	parsedToken, err := ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userEmail, ok := claims["iss"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	var user models.User
	res := db.DB.Where("email = ?", userEmail).First(&user)
	if res.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	userID := user.ID

	var answerSet struct {
		Id      int32   `json:"question_set_id"`
		Answers []string `json:"answers"`
	}

	if ctx.Bind(&answerSet) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read the body"})
		return
	}

	score := 0
	for _, answer := range answerSet.Answers {
		val, _ := strconv.Atoi(answer)
		score += val
	}

	k10 := float64(score) / float64(len(answerSet.Answers))

	// Fetch existing marked question set for the user
	var markedSet models.MarkedQuestionSet
	res = db.DB.Where("user_id = ?", userID).First(&markedSet)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// Create a new entry if none exists
		markedSet = models.MarkedQuestionSet{
			UserID: userID,
			Marked: []int32{answerSet.Id},
		}
		res = db.DB.Create(&markedSet)
		if res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
			return
		}
	} else if res.Error == nil {
		// Update existing entry
		markedSet.Marked = append(markedSet.Marked, answerSet.Id)
		res = db.DB.Save(&markedSet)
		if res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
			return
		}
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"k10": k10, "User ID": userID})
}