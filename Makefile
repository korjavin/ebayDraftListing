.PHONY: build run clean test fmt help

# Build the application
build:
	@echo "Building ebay-listing..."
	@go build -o ebay-listing cmd/main.go
	@echo "Build complete: ./ebay-listing"

# Run the application (requires photos as arguments)
run:
	@go run cmd/main.go $(ARGS)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f ebay-listing
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build  - Build the application"
	@echo "  run    - Run the application (use ARGS='photo1.jpg photo2.jpg')"
	@echo "  clean  - Remove build artifacts"
	@echo "  test   - Run tests"
	@echo "  fmt    - Format code"
	@echo "  deps   - Download and tidy dependencies"
	@echo "  help   - Show this help message"
