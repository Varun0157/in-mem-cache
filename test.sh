#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print section headers
print_header() {
    echo -e "\n${YELLOW}=== $1 ===${NC}"
}

# Function to run a command and check its status
run_test() {
    echo -e "\n${YELLOW}Running: $1${NC}"
    eval "$1"
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        exit 1
    fi
}

# Ensure we're in the project root
if [ ! -f "go.mod" ]; then
    echo -e "${RED}Error: go.mod not found. Please run this script from the project root.${NC}"
    exit 1
fi

print_header "Running Tests"

# Run go fmt to check formatting
run_test "go fmt ./..." "Code formatting check"

# Run go vet to check for common errors
run_test "go vet ./..." "Static analysis"

# Run all tests with verbose output
run_test "go test -v ./..." "Unit tests"

# Run tests with race detector
run_test "go test -race ./..." "Race condition tests"

# Run benchmarks
print_header "Running Benchmarks"
run_test "go test -bench=. -benchmem ./..." "Benchmarks"

# Run example
print_header "Running Example"
run_test "go run ./examples/main.go" "Example program"

print_header "All Tests Passed Successfully!" 