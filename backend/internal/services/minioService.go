package services

import (
	"context"
	"fmt"
	"io"
	"os"
)

type MinioService interface {
	EnsureBucket(ctx context.Context, bucket string) error
	UploadObject(ctx context.Context, bucket, objectName string, data io.Reader, size int64, contentType string) error
	DownloadObject(ctx context.Context, bucket, objectName string) (io.Reader, error)
}

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{basePath: basePath}
}

func (s *LocalStorage) EnsureBucket(ctx context.Context, bucket string) error {
	// In local storage, buckets are just directories
	path := fmt.Sprintf("%s/%s", s.basePath, bucket)
	return os.MkdirAll(path, 0755)
}

func (s *LocalStorage) UploadObject(ctx context.Context, bucket, objectName string, data io.Reader, size int64, contentType string) error {
	path := fmt.Sprintf("%s/%s/%s", s.basePath, bucket, objectName)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the data to the file
	_, err = io.Copy(file, data)
	return err
}

func (s *LocalStorage) DownloadObject(ctx context.Context, bucket, objectName string) (io.Reader, error) {
	path := fmt.Sprintf("%s/%s/%s", s.basePath, bucket, objectName)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
