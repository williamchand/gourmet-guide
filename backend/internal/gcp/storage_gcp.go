//go:build gcp

package gcp

import (
	"bytes"
	"context"
	"fmt"
	"path"

	"cloud.google.com/go/storage"
)

type CloudStorageImageStore struct {
	client     *storage.Client
	bucketName string
}

func NewCloudStorageImageStore(ctx context.Context, bucketName string) (*CloudStorageImageStore, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &CloudStorageImageStore{client: client, bucketName: bucketName}, nil
}

func (s *CloudStorageImageStore) SaveSessionImage(ctx context.Context, sessionID, fileName string, content []byte) (string, error) {
	objectPath := path.Join("session-images", sessionID, fileName)
	writer := s.client.Bucket(s.bucketName).Object(objectPath).NewWriter(ctx)
	if _, err := bytes.NewReader(content).WriteTo(writer); err != nil {
		_ = writer.Close()
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}
	return fmt.Sprintf("gs://%s/%s", s.bucketName, objectPath), nil
}
