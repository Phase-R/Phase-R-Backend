import asyncio
import ollama
from fastapi import FastAPI, Request
from fastapi.responses import StreamingResponse
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import json

app = FastAPI()

# Middleware to allow CORS for all origins
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


# Async function to stream chat responses
async def chat_stream(message: str):
    messages = [{"role": "user", "content": message}]
    response = ollama.chat(model="llama3", messages=messages, stream=True)

    # Iterating over each chunk from the response
    for chunk in response:
        content = chunk["message"]["content"]
        yield f"data: {json.dumps({'message': content})}\n\n"
        await asyncio.sleep(0.01)  # Short sleep to avoid blocking


# POST endpoint for generating diet plan
@app.post("/generate_diet")
async def get_chat_stream(request: Request):
    # Read the prompt template from the prompt.txt file
    prompt_template = ""
    with open("../configs/prompt.txt", "r") as file:
        prompt_template = file.read()

    params = await request.json()

    # Format the message with the provided parameters
    message = prompt_template.format(
        plan=params.get("plan"),
        activity=params.get("activity"),
        target_cal=params.get("target_cal"),
        target_protein=params.get("target_protein"),
        target_fat=params.get("target_fat"),
        target_carbs=params.get("target_carbs"),
        cuisine=params.get("cuisine"),
        meal_choice=params.get("meal_choice"),
    )

    return StreamingResponse(chat_stream(message), media_type="text/event-stream")


# Main entry point to run the application with multiple workers
if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "main:app", host="0.0.0.0", port=8000, workers=4
    )  # Set number of workers as needed
