package services

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	client *minio.Client
}

func NewMinioService(endpoint, accessKey, secretKey string, useSSL bool) (*MinioService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &MinioService{client: client}, nil
}

func (s *MinioService) GetObjectInfo(ctx context.Context, bucket, objectName string) (*minio.ObjectInfo, error) {
	info, err := s.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object info: %w", err)
	}

	return &info, nil
}

func (s *MinioService) DownloadObject(ctx context.Context, bucket, objectName string) (io.Reader, error) {
	object, err := s.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download object: %w", err)
	}

	return object, nil
}

func (s *MinioService) GetObject(ctx context.Context, bucket, objectName string) (*minio.Object, error) {
	object, err := s.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return object, nil
}
