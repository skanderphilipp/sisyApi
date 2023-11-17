# Makefile for managing Docker containers

# Define the service name if you want to start/stop a specific service
SERVICE_NAME=api

# Start the Docker containers
start:
	@echo "Starting Docker containers..."
	docker-compose up -d

# Stop the Docker containers
stop:
	@echo "Stopping Docker containers..."
	docker-compose stop

# Restart a specific service
restart-service:
	@echo "Restarting service $(SERVICE_NAME)..."
	docker-compose stop $(SERVICE_NAME)
	docker-compose up -d $(SERVICE_NAME)

# Stop and remove the Docker containers
down:
	@echo "Stopping and removing Docker containers..."
	docker-compose down
