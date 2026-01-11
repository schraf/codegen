all: deps vet test build

vet:
	@echo "Vetting code..."
	@go vet ./...

fmt:
	@echo "Formatting code..."
	@go fmt ./...

test:
	@echo "Running tests..."
	@go test ./...

build:
	@echo "Building..."
	@go build -o codegen ./cmd/...

run:
	@echo "Running..."
	@go run ./cmd/...

doc:
	@echo "Generating documentation..."
	@go doc -all

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

help:
	@echo "Available targets:"
	@echo "  all    - Run vet and test"
	@echo "  vet    - Vet the code"
	@echo "  fmt    - Format the code"
	@echo "  test   - Run tests"
	@echo "  build  - Build the codegen executable"
	@echo "  run    - Run the codegen application"
	@echo "  doc    - Generate documentation"
	@echo "  deps   - Install dependencies"
	@echo ""
	@echo "  help   - Show this help message"
