from fastapi import FastAPI, HTTPException
from fastapi.responses import StreamingResponse
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import httpx
import json

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],  # Restrict to trusted origins
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Pydantic model to define the expected request body structure
class DietParams(BaseModel):
    height: int
    weight: int
    bmi: float
    gender: str
    goal: str
    activity_level: str
    duration: int
    target_cal: str
    target_protein: str
    target_fat: str
    target_carbs: str
    cuisine: str
    meal_choice: str
    allergies: str

async def generate_stream(prompt: str):
    url = "http://localhost:11434/api/generate"
    async with httpx.AsyncClient() as client:
        try:
            response = await client.post(
                url,
                json={
                    "prompt": prompt,
                    "stream": True,
                    "model": "llama3",
                },
                timeout=60,
            )
            response.raise_for_status()
            async for line in response.aiter_lines():
                if line:
                    try:
                        json_response = json.loads(line)
                        if 'response' in json_response:
                            yield f"data: {json.dumps({'message': json_response['response']})}\n\n"
                    except json.JSONDecodeError as e:
                        yield f"data: {json.dumps({'error': 'JSON parsing error'})}\n\n"
        except httpx.RequestError as e:
            yield f"data: {json.dumps({'error': f'HTTP request failed: {str(e)}'})}\n\n"
        except httpx.HTTPStatusError as e:
            yield f"data: {json.dumps({'error': f'HTTP status error: {e.response.status_code}'})}\n\n"

@app.post("/generate_diet")
async def get_chat_stream(params: DietParams):
    prompt_template = (
        "Generate a meal with a {goal} plan for a person with height {height} cm, {weight} kg, "
        "BMI of {bmi}, and a weekly activity level of {activity_level}. "
        "Target calories: {target_cal} kcal over a duration of {duration} weeks. "
        "Macro Distribution: {target_protein} g protein, {target_carbs} g carbs, and {target_fat} g fat. "
        "Provide this in a JSON format over 7 days with the following structure: "
        "[{{'day1': [{{'Type': 'Breakfast', 'Time': '08:00 AM', 'Cals': 400, 'Foods': []}}]}, ...]]. "
        "Meals must include: Breakfast, Afternoon Snack, Lunch, Evening Snack, and Dinner. "
        "Use foods from {cuisine} cuisine and follow the {meal_choice} preference. "
        "Avoid mentioning anything other than the JSON."
    )

    try:
        prompt = prompt_template.format(
            goal=params.goal,
            height=params.height,
            weight=params.weight,
            bmi=params.bmi,
            activity_level=params.activity_level,
            target_cal=params.target_cal,
            target_protein=params.target_protein,
            target_carbs=params.target_carbs,
            target_fat=params.target_fat,
            duration=params.duration,
            cuisine=params.cuisine,
            meal_choice=params.meal_choice,
        )
    except KeyError as e:
        raise HTTPException(status_code=400, detail=f"Missing parameter: {str(e)}")

    return StreamingResponse(
        generate_stream(prompt),
        media_type="text/event-stream"
    )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=9000)