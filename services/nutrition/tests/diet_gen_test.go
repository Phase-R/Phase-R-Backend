package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Phase-R/Phase-R-Backend/services/nutrition/controllers"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/monthly_diet_gen", controllers.Monthly_Diet_Gen)
	router.POST("/substitute", controllers.Diet_Sub)
	return router
}

func TestGenerateDietValidParams(t *testing.T) {
	router := setupRouter()

	params := map[string]interface{}{
		"height":            180,
		"weight":            75,
		"age":               28,
		"bmi":               23.1,
		"goal":              "muscle_gain",
		"gender":            "male",
		"activity_level":    "high",
		"duration":          12,
		"target_cal":        3000,
		"target_protein":    200,
		"target_fat":        90,
		"target_carbs":      330,
		"cuisine":           "Mediterranean",
		"meal_choice":       "non-vegetarian",
		"allergies":         "none",
		"other_preferences": "low sugar, high fiber",
		"variety":           "high",
		"number_of_meals":   5,
		"meal_timings":      []string{"08:00 AM", "11:00 AM", "01:00 PM", "04:30 PM", "07:30 PM"},
	}

	body, _ := json.Marshal(params)

	req, _ := http.NewRequest("POST", "/monthly_diet_gen", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGenerateDietInvalidParams(t *testing.T) {
	router := setupRouter()

	params := map[string]interface{}{
		"height":            -1,
		"weight":            -1,
		"age":               -1,
		"bmi":               -1,
		"goal":              "",
		"gender":            "unknown",
		"activity_level":    "unknown",
		"duration":          -1,
		"target_cal":        -1,
		"target_protein":    -1,
		"target_fat":        -1,
		"target_carbs":      -1,
		"cuisine":           "unknown",
		"meal_choice":       "unknown",
		"allergies":         "none",
		"other_preferences": "none",
		"variety":           "unknown",
		"number_of_meals":   0,
		"meal_timings":      []string{},
	}

	body, _ := json.Marshal(params)

	req, _ := http.NewRequest("POST", "/monthly_diet_gen", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
}

func TestGenerateDietMissingParams(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/monthly_diet_gen", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
