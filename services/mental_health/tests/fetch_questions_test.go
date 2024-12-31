package main

import (
	// "bytes"
	// "encoding/json"
	"os"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/Phase-R/Phase-R-Backend/services/mental_health/controllers"
	"github.com/golang-jwt/jwt/v5"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/fetch_questions", controllers.FetchQuestionSet)
	return router
}

func TestFetchQuestions(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		cookie         string
		expectedStatus int
	}{
		{
			name:           "Missing Auth Cookie",
			cookie:         "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Token",
			cookie:         "nbvcdrftgyhnjmkjiuhygfvcdxsdcfv",
			expectedStatus: http.StatusUnauthorized,
		},
		// pending tests
		// {
		// 	name:           "Valid Token, User Not Found",
		// 	cookie:         createToken("nihal1anime@gmail.com"),
		// 	expectedStatus: http.StatusNotFound,
		// },
		// {
		// 	name:           "Valid Token, Successful Fetch",
		// 	cookie:         createToken("nihaltm2002@gmail.com"),
		// 	expectedStatus: http.StatusOK,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/fetch_questions", nil)
			if tt.cookie != "" {
				req.AddCookie(&http.Cookie{Name: "Auth", Value: tt.cookie})
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

// Helper function to create a JWT token for testing
func createToken(email string) string {
	claims := jwt.MapClaims{
		"iss": email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("JWT_SECRET")) // Use the same secret key as in your application
	tokenString, _ := token.SignedString(secretKey)
	return tokenString
}