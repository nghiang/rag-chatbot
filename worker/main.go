package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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
		// Update status to failed
		if err := services.UpdateDocumentStatus(payload.DocumentID, "failed"); err != nil {
			log.Printf("Failed to update document status to failed: %v", err)
		}
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
		// Update status to failed
		if err := services.UpdateDocumentStatus(payload.DocumentID, "failed"); err != nil {
			log.Printf("Failed to update document status to failed: %v", err)
		}
		return processingErr
	}

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

// processCSVFile processes CSV files by creating a schema and table in PostgreSQL
func processCSVFile(ctx context.Context, minioSvc *services.MinioService, payload ProcessDocumentPayload) error {
	log.Printf("[Worker] Processing CSV file for Knowledge Base ID: %d", payload.KnowledgeBaseID)

	// Download the CSV file from MinIO
	object, err := minioSvc.GetObject(ctx, payload.Bucket, payload.ObjectName)
	if err != nil {
		return fmt.Errorf("failed to get CSV object from MinIO: %w", err)
	}
	defer object.Close()

	// Read CSV content
	reader := csv.NewReader(object)

	// Read header row
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV headers: %w", err)
	}

	if len(headers) == 0 {
		return fmt.Errorf("CSV file has no headers")
	}

	log.Printf("[Worker] CSV Headers: %v", headers)

	// Create schema name based on knowledge base ID
	schemaName := fmt.Sprintf("kb_%d", payload.KnowledgeBaseID)

	// Create schema if it doesn't exist
	if err := services.CreateSchemaIfNotExists(schemaName); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Create table name based on document ID
	tableName := fmt.Sprintf("doc_%d", payload.DocumentID)

	// Create table with columns from CSV headers
	if err := services.CreateTableFromCSV(schemaName, tableName, headers); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Add description as table comment if provided
	if payload.Description != "" {
		if err := services.SetTableComment(schemaName, tableName, payload.Description); err != nil {
			log.Printf("[Worker] Warning: Failed to set table comment: %v", err)
			// Don't fail the entire process if comment fails
		} else {
			log.Printf("[Worker] Table comment set successfully")
		}
	}

	// Insert CSV data into the table
	rowCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV row: %w", err)
		}

		if err := services.InsertCSVRow(schemaName, tableName, headers, record); err != nil {
			log.Printf("Warning: failed to insert row: %v", err)
			continue
		}
		rowCount++
	}

	log.Printf("[Worker] CSV processing complete. Inserted %d rows into %s.%s", rowCount, schemaName, tableName)
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
