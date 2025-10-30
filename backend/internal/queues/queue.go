package queues

import (
    "context"
    "encoding/json"
    "fmt"
    "backend/config"

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

func EnqueueProcessDocument(kbID uint, documentID uint, description string, bucket, objectName string, fileType string) error {
    redisAddr := config.LoadConfig().RedisAddr

    client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
    defer client.Close()

    payload := ProcessDocumentPayload{
        KnowledgeBaseID: kbID,
        DocumentID:      documentID,
        Description:     description,
        Bucket:          bucket,
        ObjectName:      objectName,
        FileType:        fileType,
    }
    b, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    task := asynq.NewTask(TaskTypeProcessDocument, b)
    _, err = client.EnqueueContext(context.Background(), task)
    if err != nil {
        return fmt.Errorf("enqueue failed: %w", err)
    }
    return nil
}
