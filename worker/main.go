package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"worker/config"
	"worker/services"

	"github.com/hibiken/asynq"
)

const TaskTypeProcessDocument = "document:process"

type ProcessDocumentPayload struct {
	KnowledgeBaseID uint   `json:"knowledge_base_id"`
	DocumentID      uint   `json:"document_id"`
	Description     string `json:"description"`
	Bucket          string `json:"bucket"`
	ObjectName      string `json:"object_name"`
	FileType        string `json:"file_type"`
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize MinIO service
	minioSvc, err := services.NewMinioService(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioUseSSL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO service: %v", err)
	}

	log.Printf("Starting worker, connecting to Redis at %s", cfg.RedisAddr)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"default": 1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskTypeProcessDocument, func(ctx context.Context, t *asynq.Task) error {
		return handleProcessDocument(ctx, t, minioSvc)
	})

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Could not run worker server: %v", err)
	}
}

func handleProcessDocument(ctx context.Context, t *asynq.Task, minioSvc *services.MinioService) error {
	var payload ProcessDocumentPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("Failed to unmarshal payload: %v", err)
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Println("========================================")
	log.Println("[Worker] Document Processing Started")
	log.Println("========================================")
	log.Printf("Document ID:    %d", payload.DocumentID)
	log.Printf("Knowledge Base ID: %d", payload.KnowledgeBaseID)
	log.Printf("Bucket:         %s", payload.Bucket)
	log.Printf("Object Name:    %s", payload.ObjectName)
	log.Printf("File Type:     %s", payload.FileType)
	log.Printf("Task ID:        %s", t.ResultWriter().TaskID())
	log.Printf("Started At:     %s", time.Now().Format(time.RFC3339))
	log.Println("========================================")

	// Get file metadata from MinIO
	objectInfo, err := minioSvc.GetObjectInfo(ctx, payload.Bucket, payload.ObjectName)
	if err != nil {
		log.Printf("Failed to get object info: %v", err)
		return fmt.Errorf("failed to get object info: %w", err)
	}

	log.Println("[Worker] File Metadata:")
	log.Printf("  - Size:         %d bytes", objectInfo.Size)

	// Process based on file type
	var processingErr error
	switch payload.FileType {
	case "csv":
		processingErr = processCSVFile(ctx, minioSvc, payload)
	case "doc":
		processingErr = processDocFile(ctx, minioSvc, payload)
	case "graph":
		processingErr = processGraphFile(ctx, minioSvc, payload)
	default:
		processingErr = fmt.Errorf("unsupported file type: %s", payload.FileType)
	}

	if processingErr != nil {
		log.Printf("Failed to process document: %v", processingErr)
		return processingErr
	}

	// After processing is complete
	log.Println("========================================")
	log.Printf("[Worker] Document ID %d processed successfully!", payload.DocumentID)
	log.Printf("Completed At:   %s", time.Now().Format(time.RFC3339))
	log.Println("========================================")

	return nil
}

// processCSVFile processes CSV files (simulates 5s processing)
func processCSVFile(ctx context.Context, minioSvc *services.MinioService, payload ProcessDocumentPayload) error {
	log.Printf("[Worker] Processing CSV file: %s", payload.ObjectName)
	log.Println("[Worker] Simulating CSV processing... (5 seconds)")

	// Simulate CSV processing
	time.Sleep(5 * time.Second)

	log.Printf("[Worker] CSV file processing complete: %s", payload.ObjectName)
	return nil
}

// processDocFile processes document files (simulates 5s processing)
func processDocFile(ctx context.Context, minioSvc *services.MinioService, payload ProcessDocumentPayload) error {
	log.Printf("[Worker] Processing DOC file: %s", payload.ObjectName)
	log.Println("[Worker] Simulating document processing... (5 seconds)")

	// Simulate document processing
	time.Sleep(5 * time.Second)

	log.Printf("[Worker] DOC file processing complete: %s", payload.ObjectName)
	return nil
}

// processGraphFile processes graph files (simulates 5s processing)
func processGraphFile(ctx context.Context, minioSvc *services.MinioService, payload ProcessDocumentPayload) error {
	log.Printf("[Worker] Processing GRAPH file: %s", payload.ObjectName)
	log.Println("[Worker] Simulating graph processing... (5 seconds)")

	// Simulate graph processing
	time.Sleep(5 * time.Second)

	log.Printf("[Worker] GRAPH file processing complete: %s", payload.ObjectName)
	return nil
}
