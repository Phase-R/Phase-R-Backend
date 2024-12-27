package main

import (
	"bytes"
	"testing"
	"net/http"
	"encoding/json"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/Phase-R/Phase-R-Backend/services/nutrition/controllers"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/monthly_diet_gen", controllers.Monthly_Diet_Gen)
	return router
}

func TestGenerateDietValidParams(t *testing.T) {
	router := setupRouter()

	params := map[string]string{
		"plan":             "weight loss",
		"activity":         "moderate",
		"target_cal":       "2000",
		"target_protein":   "150",
		"target_fat":       "70",
		"target_carbs":     "250",
		"cuisine":          "Italian",
		"meal_choice":      "vegetarian",
		"occupation":       "office worker",
		"allergies":        "none",
		"other_preferences": "low sugar",
		"variety":          "high",
		"budget":           "medium",
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

	params := map[string]string{
		"plan":             "",
		"activity":         "unknown",
		"target_cal":       "not_a_number",
		"target_protein":   "-10",
		"target_fat":       "70",
		"target_carbs":     "250",
		"cuisine":          "Unknown",
		"meal_choice":      "unknown",
		"occupation":       "unknown",
		"allergies":        "none",
		"other_preferences": "none",
		"variety":          "unknown",
		"budget":           "unknown",
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