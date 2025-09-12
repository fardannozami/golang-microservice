#!/bin/bash

# Ensure swag is installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag CLI..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Navigate to project root
cd "$(dirname "$0")/.." || exit

# Generate Swagger documentation
echo "Generating Swagger documentation..."
swag init -g cmd/main.go -o ./docs

echo "Swagger documentation generated successfully!"
echo "Access Swagger UI at: http://localhost:8080/swagger/index.html when the service is running"