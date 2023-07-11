package minio

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type Manager struct {
	VideosBucket string
	CoverBucket  string
	client       *minio.Client
}

func NewManager(client *minio.Client, videosBucket, coverBucket string) *Manager {
	return &Manager{VideosBucket: videosBucket, CoverBucket: coverBucket, client: client}
}

func (m *Manager) UploadObject(ctx context.Context, filetype, bulkName, objectName string, reader io.Reader, size int64) error {
	var contentType string
	if filetype == "video" {
		contentType = "video/mp4"
	} else if filetype == "cover" {
		contentType = "image/jpeg"
	}
	_, err := m.client.PutObject(ctx, bulkName, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (m *Manager) GetObjectURL(ctx context.Context, bulkName, objectName string, expires time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, bulkName, objectName, expires, url.Values{})
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (m *Manager) RemoveObject(ctx context.Context, bulkName, objectName string) error {
	err := m.client.RemoveObject(ctx, bulkName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
