# Dockerfile
FROM golang:1.23

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY main.go ./

COPY . .

RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]
