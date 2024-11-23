# ###############################################################################
# # Stage 1: Build stage                                                        #
# ###############################################################################
# # Use the official Go image for building
# FROM golang:1.23 AS builder

# # Set the working directory inside the builder container
# WORKDIR /app

# # Copy dependency files first
# COPY go.mod go.sum ./

# # Download dependencies
# RUN go mod download

# # Copy the rest of the application source code
# COPY . .

# # Build the Go application as a static binary
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go


# ###############################################################################
# # Stage 2: Minimal runtime container                                          #
# ###############################################################################
# # Use a completely minimal runtime image
# FROM scratch

# # Set the working directory in the runtime container
# WORKDIR /root/

# # Copy the compiled application binary from the builder stage
# COPY --from=builder /app/main .

# # Expose the application's port
# EXPOSE 8080

# # Command to run the application
# ENTRYPOINT ["./main"]


###############################################################################
# Stage 1: Build stage                                                        #
###############################################################################
# Use the official Go image for building
FROM golang:1.23 AS builder

# Set the working directory inside the builder container
WORKDIR /app

# Copy dependency files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application as a static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go


###############################################################################
# Stage 2: Minimal runtime container                                          #
###############################################################################
# Use a completely minimal runtime image
FROM scratch

# Set the working directory in the runtime container
WORKDIR /root/

# Copy the compiled application binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file to the runtime container
COPY --from=builder /app/.env .

# Expose the application's port
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["./main"]
