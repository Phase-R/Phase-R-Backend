package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock valid diet params
func validDietParams() map[string]string {
	return map[string]string{
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
}

// Mock invalid diet params
func invalidDietParams() map[string]string {
	return map[string]string{
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
}

func TestGenerateDietValidParams(t *testing.T) {
	// Serialize valid params to JSON
	params, _ := json.Marshal(validDietParams())

	req, _ := http.NewRequest("POST", "/generate_diet", bytes.NewBuffer(params))
	req.Header.Set("Content-Type", "application/json")
	
	// Mock HTTP server
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Monthly_Diet_Gen)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGenerateDietInvalidParams(t *testing.T) {
	params, _ := json.Marshal(invalidDietParams())

	req, _ := http.NewRequest("POST", "/generate_diet", bytes.NewBuffer(params))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Monthly_Diet_Gen)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
}

func TestGenerateDietMissingParams(t *testing.T) {
	req, _ := http.NewRequest("POST", "/generate_diet", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Monthly_Diet_Gen)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
}
