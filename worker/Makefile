.PHONY: build run clean deps

# Build the worker binary
build:
	go build -o worker .

# Run the worker
run:
	go run main.go

# Clean build artifacts
clean:
	rm -f worker

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run with specific Redis address
run-dev:
	REDIS_ADDR=localhost:6379 go run main.go
