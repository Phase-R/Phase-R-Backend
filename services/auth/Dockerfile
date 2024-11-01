FROM golang:1.22.3

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY ../../db ./db

COPY main.go .

RUN go mod download
RUN go build -o bin .

ENTRYPOINT ["/app/bin"]
