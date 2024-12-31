FROM python:3.9-slim

WORKDIR /app

COPY . /app/

RUN pip install --upgrade pip && pip install -r requirements.txt

EXPOSE 9000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "9000"]