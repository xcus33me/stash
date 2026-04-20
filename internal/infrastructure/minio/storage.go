package store

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

type FileStorage interface {
	Upload(ctx context.Context, key string, r io.Reader, size int64, mimeType string) error
	Delete(ctx context.Context, key string) error
	PresignedURL(ctx context.Context, key string, ttl time.Duration) (string, error)
}

type minioStorage struct {
	client     *minio.Client
	bucketName string
}

func New(client *minio.Client, bucketName string) FileStorage {
	return &minioStorage{
		client, bucketName,
	}
}

func (s *minioStorage) Upload(ctx context.Context, key string, r io.Reader, size int64, mimeType string) error {
	_, err := s.client.PutObject(ctx, s.bucketName, key, r, size, minio.PutObjectOptions{
		ContentType: mimeType,
	})
	return err
}

func (s *minioStorage) Delete(ctx context.Context, key string) error {
	return s.client.RemoveObject(ctx, s.bucketName, key, minio.RemoveObjectOptions{})
}

func (s *minioStorage) PresignedURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	u, err := s.client.PresignedGetObject(ctx, s.bucketName, key, ttl, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
