package storage

import (
    "context"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

// Lightweight local storage fallback used when MinIO client is not available.
// This writes uploaded objects to ./uploads/<bucket>/<objectName> and is
// sufficient for local testing. To integrate MinIO, replace this file with
// an implementation that uses github.com/minio/minio-go/v7.

func EnsureBucket(ctx context.Context, bucket string) error {
    base := filepath.Join("uploads", bucket)
    return os.MkdirAll(base, 0o755)
}

func UploadObject(ctx context.Context, bucket, objectName string, reader io.Reader, size int64, contentType string) error {
    base := filepath.Join("uploads", bucket)
    if err := os.MkdirAll(base, 0o755); err != nil {
        return err
    }
    outPath := filepath.Join(base, objectName)
    if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
        return err
    }
    f, err := os.Create(outPath)
    if err != nil {
        return err
    }
    defer f.Close()
    _, err = io.Copy(f, reader)
    if err != nil {
        return fmt.Errorf("failed to write file: %w", err)
    }
    return nil
}
