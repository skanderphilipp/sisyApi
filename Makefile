# Application Name
APP_NAME=sisyApi

# Docker Image Name
DOCKER_IMAGE_NAME=myapp-image

# Build the Go app
build:
	@echo "Building the application..."
	go build -o $(APP_NAME) .

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run the application
run:
	@echo "Running the application..."
	./$(APP_NAME)

# Clean up
clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)

# Build a Docker container
docker-build:
	@echo "Building Docker container..."
	docker build -t $(DOCKER_IMAGE_NAME) .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -d --name $(APP_NAME)-container $(DOCKER_IMAGE_NAME)

# Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(APP_NAME)-container
	docker rm $(APP_NAME)-container

# Remove Docker image
docker-rmi:
	@echo "Removing Docker image..."
	docker rmi $(DOCKER_IMAGE_NAME)

# Docker: Build and run
docker-up: docker-build docker-run

# Docker: Stop and remove
docker-down: docker-stop docker-rmi

.PHONY: build test run clean docker-build docker-run docker-stop docker-rmi docker-up docker-down
