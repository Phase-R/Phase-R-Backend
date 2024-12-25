package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Phase-R/Phase-R-Backend/services/nutrition/controllers"
	"github.com/gin-gonic/gin"
)

func TestDietGenProxy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/proxy", controllers.DietGenProxy)

	t.Run("Valid Request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/proxy", nil)
		if err != nil {
			t.Fatalf("Couldn't create request: %v\n", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200 but got %d\n", w.Code)
		}
	})

	t.Run("Invalid Method", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/proxy", nil)
		if err != nil {
			t.Fatalf("Couldn't create request: %v\n", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Fatalf("Expected status 405 but got %d\n", w.Code)
		}
	})
}
