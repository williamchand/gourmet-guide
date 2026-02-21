//go:build gcp

package gcp

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/gourmet-guide/backend/internal/domain"
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

func (s *FirestoreStore) SaveSession(ctx context.Context, session domain.ConciergeSession) error {
	_, err := s.client.Collection("agent_sessions").Doc(session.ID).Set(ctx, session)
	return err
}

func (s *FirestoreStore) LoadSession(ctx context.Context, sessionID string) (domain.ConciergeSession, error) {
	snap, err := s.client.Collection("agent_sessions").Doc(sessionID).Get(ctx)
	if err != nil {
		return domain.ConciergeSession{}, err
	}
	var session domain.ConciergeSession
	if err := snap.DataTo(&session); err != nil {
		return domain.ConciergeSession{}, err
	}
	return session, nil
}

func (s *FirestoreStore) SaveMenuSafetyMetadata(ctx context.Context, restaurantID string, items []domain.MenuItem) error {
	_, _, err := s.client.Collection("menu_safety").Doc(restaurantID).Set(ctx, map[string]any{"items": items})
	return err
}

func (s *FirestoreStore) LoadMenuSafetyMetadata(ctx context.Context, restaurantID string) ([]domain.MenuItem, error) {
	snap, err := s.client.Collection("menu_safety").Doc(restaurantID).Get(ctx)
	if err != nil {
		if errors.Is(err, firestore.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	var payload struct {
		Items []domain.MenuItem `firestore:"items"`
	}
	if err := snap.DataTo(&payload); err != nil {
		return nil, err
	}
	return payload.Items, nil
}

func (s *FirestoreStore) SaveImageReference(ctx context.Context, sessionID, imagePath string) error {
	_, _, err := s.client.Collection("agent_sessions").Doc(sessionID).Set(ctx, map[string]any{
		"imageRefs": firestore.ArrayUnion(imagePath),
	}, firestore.MergeAll)
	return err
}

func (s *FirestoreStore) Close() error { return s.client.Close() }
