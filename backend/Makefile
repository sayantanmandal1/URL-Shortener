# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=url-shortener
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
test-race:
	$(GOTEST) -v -race ./...

# Run benchmarks
bench:
	$(GOTEST) -v -bench=. ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./
	./$(BINARY_NAME)

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Docker build
docker-build:
	docker build -t $(BINARY_NAME) .

# Docker run
docker-run:
	docker run -p 8000:8000 $(BINARY_NAME)

# Lint code
lint:
	golangci-lint run

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Security check
security:
	gosec ./...

# All checks (format, vet, lint, test)
check: fmt vet test

# CI pipeline
ci: deps fmt vet test-race test-coverage

.PHONY: build clean test test-coverage test-race bench deps run build-linux docker-build docker-run lint fmt vet security check ci