###############################################################################
# Stage 1: Build stage                                                        #
###############################################################################
# Use the official Go image for building
FROM golang:1.23 AS builder

# Set the working directory inside the builder container
WORKDIR /app

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application for Linux
RUN GOOS=linux GOARCH=amd64 go build -o main main.go


###############################################################################
# Stage 2: Minimal runtime container                                          #
###############################################################################
# Use a minimal runtime image
FROM debian:bullseye-slim

# Set the working directory in the runtime container
WORKDIR /root/

# Copy the compiled application binary from the builder stage
COPY --from=builder /app/main .

# Expose the application's port
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["./main"]
