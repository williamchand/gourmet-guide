//go:build !gcp

package seed

import (
	"context"
	"errors"

	"github.com/gourmet-guide/backend/internal/domain"
)

var errGCPBuildRequired = errors.New("gcp build tag required for GCS/Firestore publishing")

// GCSImageUploader is unavailable without gcp build tag.
type GCSImageUploader struct{}

func NewGCSImageUploader(_ context.Context, _ string) (*GCSImageUploader, error) {
	return nil, errGCPBuildRequired
}

func (u *GCSImageUploader) UploadMenuImage(_ context.Context, _, _, _ string, _ []byte) (string, error) {
	return "", errGCPBuildRequired
}

func (u *GCSImageUploader) Close() error { return nil }

// FirestoreRestaurantStore is unavailable without gcp build tag.
type FirestoreRestaurantStore struct{}

func NewFirestoreRestaurantStore(_ context.Context, _, _ string) (*FirestoreRestaurantStore, error) {
	return nil, errGCPBuildRequired
}

func (s *FirestoreRestaurantStore) SaveRestaurant(_ context.Context, _ domain.Restaurant) error {
	return errGCPBuildRequired
}

func (s *FirestoreRestaurantStore) Close() error { return nil }
