from fastapi import FastAPI, Request
from fastapi.responses import StreamingResponse
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import requests
import json

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Pydantic model to define the expected request body structure
class DietParams(BaseModel):
    plan: str
    activity: str
    target_cal: str
    target_protein: str
    target_fat: str
    target_carbs: str
    cuisine: str
    meal_choice: str

async def generate_stream(prompt: str):
    # Make a streaming request to Ollama
    response = requests.post(
        "http://localhost:11434/api/generate", 
        json={
            "prompt": prompt,
            "stream": True,
            "model": "llama3",
        }, 
        stream=True
    )

    # Stream the response
    for line in response.iter_lines():
        if line:
            try:
                # Decode and parse the JSON line
                json_response = json.loads(line.decode('utf-8'))
                
                # Check if the response contains content
                if 'response' in json_response:
                    # Yield the content as a server-sent event
                    yield f"data: {json.dumps({'message': json_response['response']})}\n\n"
            except Exception as e:
                # Handle any parsing errors
                yield f"data: {json.dumps({'error': str(e)})}\n\n"

@app.post("/generate_diet")
async def get_chat_stream(params: DietParams):
    prompt_template = "Generate a meal plan for {plan} with a daily activity level being {activity}. Target calories: {target_cal} kcal. Macro Distribution: {target_protein} g protein, {target_carbs} g of carbs and {target_fat} g of fat. Give this in the form a table with days: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday] with meals being mandatorily [Breakfast, Afternoon Snack, Lunch, Evening Snack, Dinner]. Don't tell anything other than the table in the form of html table as mentioned. The foods should mainly belong to {cuisine} cuisine and should be {meal_choice}."
    
    prompt = prompt_template.format(
        plan=params.plan,
        activity=params.activity,
        target_cal=params.target_cal,
        target_protein=params.target_protein,
        target_fat=params.target_fat,
        target_carbs=params.target_carbs,
        cuisine=params.cuisine,
        meal_choice=params.meal_choice
    )
    
    return StreamingResponse(
        generate_stream(prompt), 
        media_type="text/event-stream"
    )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=9000)