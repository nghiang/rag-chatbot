package queues

import (
    "context"
    "encoding/json"
    "fmt"
    "os"

    "github.com/hibiken/asynq"
)

const TaskTypeProcessDocument = "document:process"

type ProcessDocumentPayload struct {
    DocumentID uint   `json:"document_id"`
    Bucket     string `json:"bucket"`
    ObjectName string `json:"object_name"`
}

func EnqueueProcessDocument(documentID uint, bucket, objectName string) error {
    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "127.0.0.1:6379"
    }

    client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
    defer client.Close()

    payload := ProcessDocumentPayload{
        DocumentID: documentID,
        Bucket:     bucket,
        ObjectName: objectName,
    }
    b, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    task := asynq.NewTask(TaskTypeProcessDocument, b)
    // enqueue immediately
    _, err = client.EnqueueContext(context.Background(), task)
    if err != nil {
        return fmt.Errorf("enqueue failed: %w", err)
    }
    return nil
}
