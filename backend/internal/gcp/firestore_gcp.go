//go:build gcp

package gcp

import (
	"context"

	"cloud.google.com/go/firestore"
)

// FirestoreStore uses Firestore (managed GCP service) for session persistence.
type FirestoreStore struct {
	client *firestore.Client
}

func NewFirestoreStore(ctx context.Context, projectID string) (*FirestoreStore, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &FirestoreStore{client: client}, nil
}

func (s *FirestoreStore) SavePrompt(ctx context.Context, sessionID, prompt string) error {
	_, _, err := s.client.Collection("agent_sessions").Doc(sessionID).Set(ctx, map[string]any{"lastPrompt": prompt}, firestore.MergeAll)
	return err
}

func (s *FirestoreStore) Close() error { return s.client.Close() }
