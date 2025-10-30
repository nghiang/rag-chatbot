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
	DocumentID uint   `json:"document_id"`
	Bucket     string `json:"bucket"`
	ObjectName string `json:"object_name"`
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize PostgreSQL connection
	if err := services.InitDB(
		cfg.PostGresUser,
		cfg.PostGresPassword,
		cfg.PostGresDB,
		cfg.PostGresHost,
		cfg.PostGresPort,
	); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer services.CloseDB()

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
	log.Printf("Bucket:         %s", payload.Bucket)
	log.Printf("Object Name:    %s", payload.ObjectName)
	log.Printf("Task ID:        %s", t.ResultWriter().TaskID())
	log.Printf("Started At:     %s", time.Now().Format(time.RFC3339))
	log.Println("========================================")

	// Get file metadata from MinIO
	objectInfo, err := minioSvc.GetObjectInfo(ctx, payload.Bucket, payload.ObjectName)
	if err != nil {
		log.Printf("Failed to get object info: %v", err)
		// Update status to failed
		if err := services.UpdateDocumentStatus(payload.DocumentID, "failed"); err != nil {
			log.Printf("Failed to update document status to failed: %v", err)
		}
		return fmt.Errorf("failed to get object info: %w", err)
	}

	log.Println("[Worker] File Metadata:")
	log.Printf("  - Size:         %d bytes", objectInfo.Size)
	log.Printf("  - Content-Type: %s", objectInfo.ContentType)
	log.Printf("  - ETag:         %s", objectInfo.ETag)
	log.Printf("  - Last-Modified: %s", objectInfo.LastModified.Format(time.RFC3339))
	log.Println("========================================")

	// Simulate document processing (sleep for 5 seconds)
	log.Println("[Worker] Processing document... (simulating 5 seconds)")
	time.Sleep(5 * time.Second)

	// Update document status to processed in PostgreSQL
	if err := services.UpdateDocumentStatus(payload.DocumentID, "processed"); err != nil {
		log.Printf("Failed to update document status: %v", err)
		return fmt.Errorf("failed to update document status: %w", err)
	}

	// After processing is complete
	log.Println("========================================")
	log.Printf("[Worker] Document ID %d processed successfully!", payload.DocumentID)
	log.Printf("Completed At:   %s", time.Now().Format(time.RFC3339))
	log.Println("========================================")

	return nil
}
