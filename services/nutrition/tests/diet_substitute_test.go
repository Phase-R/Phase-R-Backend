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

func setUpRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/substitute", controllers.Substitute)
	return router
}

func TestSubstituteValidParams(t *testing.T) {
	router := setUpRouter()

	params := map[string]string{
		"food":              "chicken tikka masala",
		"allergies":         "mushrooms",
		"other_preferences": "high spice",
	}

	body, _ := json.Marshal(params)

	req, _ := http.NewRequest("POST", "/substitute", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestSubstituteInvalidParams(t *testing.T) {
	router := setUpRouter()

	params := map[string]string{
		"food":              "",
		"allergies":         "unknown",
		"other_preferences": "unknown",
	}

	body, _ := json.Marshal(params)

	req, _ := http.NewRequest("POST", "/substitute", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestSubstituteMissingParams(t *testing.T) {
	router := setUpRouter()

	params := map[string]string{
		"food":              "chicken tikka masala",
		"allergies":         "mushrooms",
	}

	body, _ := json.Marshal(params)

	req, _ := http.NewRequest("POST", "/substitute", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}