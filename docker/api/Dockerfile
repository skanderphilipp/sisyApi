# Start from the official Go image to create a build artifact.
FROM golang:1.21.4-alpine AS base

WORKDIR /app

ENV GO111MODULE="on"
ENV GOOS="linux"
ENV CGO_ENABLED=0
# Set the environment variable to ensure binaries are in the PATH
# System dependencies
RUN apk update \
  && apk add --no-cache \
  ca-certificates \
  git \
  && update-ca-certificates

### Development with hot reload and debugger
FROM base AS dev
ENV PATH="/go/bin:${PATH}"
WORKDIR /app
# Copy the Go modules and sum files
COPY go.mod go.sum ./
# Download Go module dependencies
RUN go mod download

# Install air for live reloading during development
RUN go install github.com/cosmtrek/air@latest && which air

# Install delve for debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN chmod +x /go/bin/air
EXPOSE 8080
EXPOSE 2345

ENTRYPOINT ["air"]

### Executable builder
FROM base AS builder
WORKDIR /app

# Application dependencies
COPY . /app
RUN go mod download \
  && go mod verify

RUN go build -o my-great-program -a .

### Production
FROM alpine:latest

RUN apk update \
  && apk add --no-cache \
  ca-certificates \
  curl \
  tzdata \
  && update-ca-certificates

# Copy executable
COPY --from=builder /app/my-great-program /usr/local/bin/my-great-program
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/my-great-program"]