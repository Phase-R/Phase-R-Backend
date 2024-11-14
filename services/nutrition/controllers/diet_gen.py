import asyncio
import ollama
from fastapi import FastAPI, Request
from fastapi.responses import StreamingResponse
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import json
from concurrent.futures import ThreadPoolExecutor
from typing import List, Dict
import asyncio
from functools import partial

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Configure thread pool for parallel processing
THREAD_POOL_SIZE = 4
thread_pool = ThreadPoolExecutor(max_workers=THREAD_POOL_SIZE)


class DietParams(BaseModel):
    plan: str
    activity: str
    target_cal: str
    target_protein: str
    target_fat: str
    target_carbs: str
    cuisine: str
    meal_choice: str


async def process_ollama_response(response) -> str:
    """Process individual Ollama response chunks"""
    result = ""
    try:
        for chunk in response:
            content = chunk["message"]["content"]
            result += content
    except Exception as e:
        print(f"Error processing response: {e}")
    return result


async def generate_diet_parallel(message: str, batch_size: int = 3) -> str:
    """Generate diet suggestions in parallel using multiple Ollama instances"""
    messages = [{"role": "user", "content": message}]

    # Create multiple concurrent requests to Ollama
    async def make_single_request():
        response = ollama.chat(model="llama3", messages=messages)
        return response["message"]["content"]

    # Create multiple tasks
    tasks = [make_single_request() for _ in range(batch_size)]

    # Wait for all tasks to complete
    responses = await asyncio.gather(*tasks, return_exceptions=True)

    # Filter out any failed responses and combine results
    valid_responses = [r for r in responses if isinstance(r, str)]
    if not valid_responses:
        return "Failed to generate diet suggestions"

    # Return the best or most complete response
    return max(valid_responses, key=len)


async def chat_stream(message: str):
    """Stream the diet generation results"""
    try:
        # Generate diet suggestions in parallel
        result = await generate_diet_parallel(message)

        # Stream the result in chunks
        chunk_size = 100
        for i in range(0, len(result), chunk_size):
            chunk = result[i : i + chunk_size]
            yield f"data: {json.dumps({'message': chunk})}\n\n"
            await asyncio.sleep(0.01)

    except Exception as e:
        yield f"data: {json.dumps({'error': str(e)})}\n\n"


class ResponseCache:
    def __init__(self):
        self.cache = {}
        self.lock = asyncio.Lock()

    async def get(self, key: str) -> str:
        async with self.lock:
            return self.cache.get(key)

    async def set(self, key: str, value: str):
        async with self.lock:
            self.cache[key] = value


# Initialize cache
response_cache = ResponseCache()


@app.post("/generate_diet")
async def get_chat_stream(request: Request):
    try:
        # Read the prompt template
        prompt_template = ""
        with open("../configs/prompt.txt", "r") as file:
            prompt_template = file.read()

        # Get parameters from request
        params = await request.json()

        # Create cache key based on parameters
        cache_key = json.dumps(params, sort_keys=True)

        # Check cache first
        cached_response = await response_cache.get(cache_key)
        if cached_response:

            async def stream_cached():
                yield f"data: {json.dumps({'message': cached_response})}\n\n"

            return StreamingResponse(stream_cached(), media_type="text/event-stream")

        # Format message with parameters
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

    except Exception as e:
        return StreamingResponse(
            iter([f"data: {json.dumps({'error': str(e)})}\n\n"]),
            media_type="text/event-stream",
        )


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        app,
        host="0.0.0.0",
        port=8000,
        workers=4,  # Multiple worker processes for handling concurrent requests
    )
