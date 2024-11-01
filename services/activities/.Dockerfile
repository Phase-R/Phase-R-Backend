# Dockerfile
FROM golang:1.23

# Define the build argument
ARG DB_PATH=../../db

# Set the working directory
WORKDIR /app

# Copy the main application files
COPY go.mod go.sum ./
COPY main.go ./

# Copy the db module using the build argument
COPY ${DB_PATH} ./db

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o bin .

# Specify the entry point
ENTRYPOINT ["/app/bin"]
