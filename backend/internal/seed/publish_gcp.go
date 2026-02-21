//go:build gcp

package seed

import (
	"context"
	"fmt"
	"path"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/gourmet-guide/backend/internal/domain"
)

// GCSImageUploader uploads image artifacts to GCS.
type GCSImageUploader struct {
	bucket string
	client *storage.Client
}

func NewGCSImageUploader(ctx context.Context, bucket string) (*GCSImageUploader, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GCSImageUploader{bucket: bucket, client: client}, nil
}

func (u *GCSImageUploader) UploadMenuImage(ctx context.Context, restaurantID, menuItemID, fileName string, data []byte) (string, error) {
	objectPath := path.Join("seed-images", restaurantID, menuItemID, fileName)
	writer := u.client.Bucket(u.bucket).Object(objectPath).NewWriter(ctx)
	if _, err := writer.Write(data); err != nil {
		_ = writer.Close()
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}
	return fmt.Sprintf("gs://%s/%s", u.bucket, objectPath), nil
}

func (u *GCSImageUploader) Close() error {
	return u.client.Close()
}

// FirestoreRestaurantStore writes generated restaurants to Firestore.
type FirestoreRestaurantStore struct {
	collection string
	client     *firestore.Client
}

func NewFirestoreRestaurantStore(ctx context.Context, projectID, collection string) (*FirestoreRestaurantStore, error) {
	if collection == "" {
		collection = "restaurants"
	}
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &FirestoreRestaurantStore{collection: collection, client: client}, nil
}

func (s *FirestoreRestaurantStore) SaveRestaurant(ctx context.Context, restaurant domain.Restaurant) error {
	_, err := s.client.Collection(s.collection).Doc(restaurant.ID).Set(ctx, restaurant)
	return err
}

func (s *FirestoreRestaurantStore) Close() error {
	return s.client.Close()
}
