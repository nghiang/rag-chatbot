# Document Processing Worker

This is a standalone worker service that processes document embedding tasks asynchronously using Redis and Asynq.

## Features

- Processes document embedding tasks from a Redis queue
- Simulates document processing with a 5-second delay
- Logs processing status and completion messages
- Can be run independently from the backend service

## Prerequisites

- Go 1.21 or higher
- Redis server running

## Setup

1. Install dependencies:
```bash
make deps
```

2. Build the worker:
```bash
make build
```

## Running

### Development Mode

```bash
make run-dev
```

Or with custom Redis address:
```bash
REDIS_ADDR=localhost:6379 go run main.go
```

### Production Mode

```bash
REDIS_ADDR=redis:6379 ./worker
```

### Using Docker

Build the Docker image:
```bash
docker build -t document-worker .
```

Run the container:
```bash
docker run -e REDIS_ADDR=redis:6379 document-worker
```

## Configuration

The worker can be configured using environment variables:

- `REDIS_ADDR`: Redis server address (default: `127.0.0.1:6379`)

## Task Processing

The worker listens for tasks of type `document:process` with the following payload structure:

```json
{
  "document_id": 123,
  "bucket": "documents",
  "object_name": "kb_1/12345_file.pdf"
}
```

When a task is received:
1. The worker logs the start of processing
2. Simulates processing for 5 seconds
3. Logs completion with a success message
4. Returns the result to the queue

## Integration with Backend

This worker integrates with the backend service through the Redis queue. The backend enqueues document processing tasks using the `queues.EnqueueProcessDocument()` function, and this worker picks them up for processing.
