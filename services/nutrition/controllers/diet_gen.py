import asyncio
import json

import ollama
from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import StreamingResponse
from pydantic import BaseModel

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
    occupation: str
    allergies: str
    other_preferences: str
    variety: str
    budget: str


async def chat_stream(message: str):
    messages = [{"role": "user", "content": message}]
    response = ollama.chat(model='llama3', messages=messages, stream=True)

    for chunk in response:
        content = chunk['message']['content']
        yield f"data: {json.dumps({'message': content})}\n\n"
        await asyncio.sleep(0.01)


@app.post("/generate_diet")
async def get_chat_stream(request: Request):
    # Read the prompt template from the prompt.txt file
    # Update: The prompt.txt file is no longer used.
    # The prompt is called as a string in the program itself.
    # Only use the text file for testing purposes.
    prompt_template = (
        f"""Generate a meal plan for {{plan}} with a daily activity level being {{activity}}.
        Target calories: {{target_cal}} kcal. Macro Distribution: {{target_protein}} g protein, 
        {{target_carbs}} g of carbs and {{target_fat}} g of fat. Give this in the form a table with days: 
        [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday] with meals being mandatory 
        [Breakfast, Afternoon Snack, Lunch, Evening Snack, Dinner]. Don't tell anything other than the table 
        in the form of html table as mentioned. The foods should mainly belong to {{cuisine}} cuisine and should be {{meal_choice}},
        keeping in mind the occupation of the person as an {{occupation}}. The person has {{allergies}} allergies and
        other preferences include {{other_preferences}}. The person prefers a {{variety}} variety of foods and has a
        {{budget}} budget. The meal plan should be for a week."""
    )
    # with open("../configs/prompt.txt", "r") as file:
    #     prompt_template = file.read()

    params = await request.json()

    message = prompt_template.format(
        plan=params.get("plan"),
        activity=params.get("activity"),
        target_cal=params.get("target_cal"),
        target_protein=params.get("target_protein"),
        target_fat=params.get("target_fat"),
        target_carbs=params.get("target_carbs"),
        cuisine=params.get("cuisine"),
        meal_choice=params.get("meal_choice"),
        occupation=params.get("occupation"),
        allergies=params.get("allergies"),
        other_preferences=params.get("other_preferences"),
        variety=params.get("variety"),
        budget=params.get("budget")
    )

    return StreamingResponse(chat_stream(message), media_type="text/event-stream")


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
