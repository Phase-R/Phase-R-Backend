import uvicorn

if __name__ == "__main__":
    uvicorn.run(
        "diet_gen:app",
        host="0.0.0.0",
        port=8000,
        workers=4,   # only convert it to true for dev, or else will
        reload=False  # ignore workers flag, thereby running in a single thread
    )
