.PHONY: build run clean test fmt lint vet check deps help

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

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run linters
lint:
	@echo "Checking code formatting..."
	@UNFORMATTED=$$(gofmt -s -l .); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "The following files need formatting:"; \
		echo "$$UNFORMATTED"; \
		echo "Run 'make fmt' to fix"; \
		exit 1; \
	fi
	@echo "Running go vet..."
	@go vet ./...
	@if command -v staticcheck >/dev/null 2>&1; then \
		echo "Running staticcheck..."; \
		staticcheck ./...; \
	else \
		echo "staticcheck not installed, skipping (install with: go install honnef.co/go/tools/cmd/staticcheck@latest)"; \
	fi

# Run all checks (format, lint, test, build)
check: lint test build
	@echo "All checks passed!"

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
	@echo "  vet    - Run go vet"
	@echo "  lint   - Run all linters (format check, vet, staticcheck)"
	@echo "  check  - Run all checks (lint, test, build)"
	@echo "  deps   - Download and tidy dependencies"
	@echo "  help   - Show this help message"
