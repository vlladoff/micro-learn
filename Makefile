.PHONY: help build run test clean docker-build docker-run docker-stop docker-clean

  # Default target
  help:
        @echo "Available targets:"
        @echo "  build        - Build the application"
        @echo "  run          - Run the application locally"
        @echo "  test         - Run tests"
        @echo "  clean        - Clean build artifacts"
        @echo "  docker-build - Build Docker image"
        @echo "  docker-run   - Run with Docker Compose"
        @echo "  docker-stop  - Stop Docker Compose"
        @echo "  docker-clean - Clean Docker resources"

  # Go targets
  build:
        @echo "Building application..."
        @go build -o bin/app ./cmd/smpl-api-oapi/

  run:
        @echo "Running application..."
        @go run ./cmd/smpl-api-oapi/

  test:
        @echo "Running tests..."
        @go test -v ./...

  clean:
        @echo "Cleaning build artifacts..."
        @rm -rf bin/
        @go clean

  # Docker targets
  docker-build:
        @echo "Building Docker image..."
        @docker build -t micro-learn:latest .

  docker-run:
        @echo "Starting services with Docker Compose..."
        @docker-compose up -d

  docker-stop:
        @echo "Stopping Docker Compose services..."
        @docker-compose down

  docker-clean:
        @echo "Cleaning Docker resources..."
        @docker-compose down -v --remove-orphans
        @docker system prune -f