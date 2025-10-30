# Document Processing Worker

This is a standalone worker service that processes document embedding tasks asynchronously using Redis (Asynq), MinIO, and PostgreSQL.

## Features

- Processes document embedding tasks from a Redis queue
- Fetches file metadata from MinIO
- Simulates document processing with a 5-second delay
- Updates document status in PostgreSQL database
- Structured code with separate config, MinIO, and PostgreSQL services
- Logs processing status and completion messages

## Project Structure

```
worker/
├── main.go                          # Main application entry point
├── config/
│   └── config.go                    # Configuration management
├── services/
│   ├── minioService.go              # MinIO service for object storage
│   └── postgresConnection.go       # PostgreSQL database connection
├── go.mod                           # Go module dependencies
├── Dockerfile                       # Docker container definition
├── Makefile                         # Build and run commands
└── README.md                        # This file
```

## Prerequisites

- Go 1.23 or higher
- Redis server running
- PostgreSQL database running
- MinIO server running

## Setup

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Update the `.env` file with your configuration

3. Install dependencies:
```bash
make deps
```

4. Build the worker:
```bash
make build
```

## Running

### Development Mode

```bash
make run-dev
```

Or with custom environment variables:
```bash
REDIS_ADDR=localhost:6379 \
POSTGRES_HOST=localhost \
POSTGRES_PORT=5432 \
MINIO_ENDPOINT=localhost:9008 \
go run main.go
```

### Production Mode

```bash
./worker
```

### Using Docker

Build the Docker image:
```bash
docker build -t document-worker .
```

Run the container:
```bash
docker run --env-file .env document-worker
```

## Configuration

The worker can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `REDIS_ADDR` | Redis server address | `127.0.0.1:6379` |
| `POSTGRES_USER` | PostgreSQL username | `postgres` |
| `POSTGRES_PASSWORD` | PostgreSQL password | `postgres` |
| `POSTGRES_DB` | PostgreSQL database name | `rag_chatbot_db` |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `MINIO_ENDPOINT` | MinIO server endpoint | `localhost:9008` |
| `MINIO_ACCESS_KEY` | MinIO access key | `minioadmin` |
| `MINIO_SECRET_KEY` | MinIO secret key | `minioadmin` |
| `MINIO_USE_SSL` | Use SSL for MinIO | `false` |

## Task Processing

The worker listens for tasks of type `document:process` with the following payload structure:

```json
{
  "document_id": 123,
  "bucket": "documents",
  "object_name": "kb_1/12345_file.pdf"
}
```

### Processing Flow

1. Worker receives task from Redis queue
2. Fetches file metadata from MinIO (size, content-type, etc.)
3. Logs file information
4. Simulates processing for 5 seconds
5. Updates document status to "processed" in PostgreSQL
6. Logs completion

### Example Output

```
========================================
[Worker] Document Processing Started
========================================
Document ID:    123
Bucket:         documents
Object Name:    kb_1/12345_file.pdf
Task ID:        abc-123-xyz
Started At:     2025-10-30T12:00:00Z
========================================
[Worker] File Metadata:
  - Size:         1048576 bytes
  - Content-Type: application/pdf
  - ETag:         "d41d8cd98f00b204e9800998ecf8427e"
  - Last-Modified: 2025-10-30T11:55:00Z
========================================
[Worker] Processing document... (simulating 5 seconds)
Updated document ID 123 status to: processed
========================================
[Worker] Document ID 123 processed successfully!
Completed At:   2025-10-30T12:00:05Z
========================================
```

## Integration with Backend

This worker integrates with the backend service through:
- **Redis**: Backend enqueues tasks using `queues.EnqueueProcessDocument()`
- **MinIO**: Shared object storage for uploaded files
- **PostgreSQL**: Shared database for document metadata

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
make build
```

### Cleaning Build Artifacts
```bash
make clean
```

