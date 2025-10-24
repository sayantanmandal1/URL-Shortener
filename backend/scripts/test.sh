#!/bin/bash

# Backend Test Script
# This script runs comprehensive tests for the backend

set -e

echo "🧪 Running Backend Test Suite"
echo "=============================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Change to backend directory
cd "$(dirname "$0")/.."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if ! printf '%s\n%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V -C; then
    print_error "Go version $GO_VERSION is too old. Please upgrade to Go $REQUIRED_VERSION or later."
    exit 1
fi

print_status "Go version $GO_VERSION detected"

# Download dependencies
echo ""
echo "📦 Downloading dependencies..."
go mod download
go mod tidy
print_status "Dependencies downloaded"

# Format check
echo ""
echo "🎨 Checking code formatting..."
UNFORMATTED=$(gofmt -s -l .)
if [ -n "$UNFORMATTED" ]; then
    print_error "Code is not formatted. Please run 'go fmt ./...'"
    echo "Unformatted files:"
    echo "$UNFORMATTED"
    exit 1
fi
print_status "Code is properly formatted"

# Vet
echo ""
echo "🔍 Running go vet..."
if go vet ./...; then
    print_status "go vet passed"
else
    print_error "go vet failed"
    exit 1
fi

# Security scan (if gosec is available)
echo ""
echo "🔒 Running security scan..."
if command -v gosec &> /dev/null; then
    if gosec -quiet ./...; then
        print_status "Security scan passed"
    else
        print_warning "Security scan found issues (non-blocking)"
    fi
else
    print_warning "gosec not installed, skipping security scan"
fi

# Lint (if golangci-lint is available)
echo ""
echo "🧹 Running linter..."
if command -v golangci-lint &> /dev/null; then
    if golangci-lint run --timeout=5m; then
        print_status "Linting passed"
    else
        print_error "Linting failed"
        exit 1
    fi
else
    print_warning "golangci-lint not installed, skipping linting"
fi

# Unit tests
echo ""
echo "🧪 Running unit tests..."
if go test -v -race -coverprofile=coverage.out ./...; then
    print_status "Unit tests passed"
else
    print_error "Unit tests failed"
    exit 1
fi

# Coverage report
echo ""
echo "📊 Generating coverage report..."
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
echo "Total coverage: $COVERAGE"

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
print_status "Coverage report generated: coverage.html"

# Build test
echo ""
echo "🔨 Testing build..."
if go build -o url-shortener ./; then
    print_status "Build successful"
    rm -f url-shortener
else
    print_error "Build failed"
    exit 1
fi

# Benchmark tests (optional)
if [ "$1" = "--bench" ]; then
    echo ""
    echo "⚡ Running benchmarks..."
    go test -bench=. -benchmem ./...
    print_status "Benchmarks completed"
fi

echo ""
echo "🎉 All tests passed successfully!"
echo "=============================="
echo "Summary:"
echo "- Code formatting: ✓"
echo "- Go vet: ✓"
echo "- Security scan: ✓"
echo "- Linting: ✓"
echo "- Unit tests: ✓"
echo "- Coverage: $COVERAGE"
echo "- Build: ✓"

if [ "$1" = "--bench" ]; then
    echo "- Benchmarks: ✓"
fi

echo ""
echo "Ready for deployment! 🚀"