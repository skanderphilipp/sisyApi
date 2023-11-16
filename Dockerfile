# Start from the official Go image to create a build artifact.
FROM golang:1.21.4-alpine as builder

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum files (if using Go modules).
COPY go.mod go.sum ./

# Download dependencies (if using Go modules).
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container.
COPY . .

# Build the Go app. Adjust the path to where main.go is located.
RUN go build -o sisyApp ./cmd

# Start a new stage from a smaller image.
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /app/sisyApp .

# Command to run the executable.
CMD ["./sisyApp"]
