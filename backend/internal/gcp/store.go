package gcp

import "context"

// SessionStore persists agent session metadata.
type SessionStore interface {
	SavePrompt(ctx context.Context, sessionID, prompt string) error
	Close() error
}

// MemoryStore is local default storage for development and tests.
type MemoryStore struct{}

func NewMemoryStore() *MemoryStore { return &MemoryStore{} }

func (m *MemoryStore) SavePrompt(_ context.Context, _, _ string) error { return nil }

func (m *MemoryStore) Close() error { return nil }
