import pytest
from fastapi.testclient import TestClient

from services.nutrition.controllers.diet_gen import app

client = TestClient(app)


@pytest.fixture
def valid_diet_params():
    return {
        "plan": "weight loss",
        "activity": "moderate",
        "target_cal": "2000",
        "target_protein": "150",
        "target_fat": "70",
        "target_carbs": "250",
        "cuisine": "Italian",
        "meal_choice": "vegetarian",
        "occupation": "office worker",
        "allergies": "none",
        "other_preferences": "low sugar",
        "variety": "high",
        "budget": "medium"
    }


@pytest.fixture
def invalid_diet_params():
    return {
        "plan": "",
        "activity": "unknown",
        "target_cal": "not_a_number",
        "target_protein": "-10",
        "target_fat": "70",
        "target_carbs": "250",
        "cuisine": "Unknown",
        "meal_choice": "unknown",
        "occupation": "unknown",
        "allergies": "none",
        "other_preferences": "none",
        "variety": "unknown",
        "budget": "unknown"
    }


def test_generate_diet_valid_params(valid_diet_params):
    response = client.post("/generate_diet", json=valid_diet_params)
    assert response.status_code == 200
    assert response.headers["content-type"] == "text/event-stream"


def test_generate_diet_invalid_params(invalid_diet_params):
    response = client.post("/generate_diet", json=invalid_diet_params)
    assert response.status_code == 422  # Unprocessable Entity


def test_generate_diet_missing_params():
    response = client.post("/generate_diet", json={})
    assert response.status_code == 422  # Unprocessable Entity
