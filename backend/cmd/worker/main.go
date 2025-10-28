package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/hibiken/asynq"

    "backend/internal/database"
    "backend/internal/models"
    "backend/internal/queues"
)

func main() {
    if err := database.InitDB(); err != nil {
        panic(err)
    }
    defer database.CloseDB()

    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "127.0.0.1:6379"
    }

    srv := asynq.NewServer(
        asynq.RedisClientOpt{Addr: redisAddr},
        asynq.Config{
            Concurrency: 10,
        },
    )

    mux := asynq.NewServeMux()
    mux.HandleFunc(queues.TaskTypeProcessDocument, func(ctx context.Context, t *asynq.Task) error {
        var p queues.ProcessDocumentPayload
        if err := json.Unmarshal(t.Payload(), &p); err != nil {
            return err
        }
        fmt.Printf("[worker] processing document id=%d object=%s\n", p.DocumentID, p.ObjectName)

        // fake processing: set status processing (already set), sleep 5s, then set done
        time.Sleep(5 * time.Second)

        if err := models.UpdateDocumentEmbeddingStatus(p.DocumentID, "done"); err != nil {
            return err
        }
        fmt.Printf("[worker] done document id=%d\n", p.DocumentID)
        return nil
    })

    if err := srv.Run(mux); err != nil {
        fmt.Fprintf(os.Stderr, "asynq server failed: %v\n", err)
        os.Exit(1)
    }
}
