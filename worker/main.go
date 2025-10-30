package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hibiken/asynq"
)

const TaskTypeProcessDocument = "document:process"

type ProcessDocumentPayload struct {
	DocumentID uint   `json:"document_id"`
	Bucket     string `json:"bucket"`
	ObjectName string `json:"object_name"`
}

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	log.Printf("Starting worker, connecting to Redis at %s", redisAddr)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"default": 1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskTypeProcessDocument, handleProcessDocument)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("Could not run worker server: %v", err)
	}
}

func handleProcessDocument(ctx context.Context, t *asynq.Task) error {
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

	// Simulate document processing (sleep for 5 seconds)
	log.Println("[Worker] Processing document... (simulating 5 seconds)")
	time.Sleep(5 * time.Second)

	// After processing is complete
	log.Println("========================================")
	log.Printf("[Worker] Document ID %d processed successfully!", payload.DocumentID)
	log.Printf("Completed At:   %s", time.Now().Format(time.RFC3339))
	log.Println("========================================")

	return nil
}
