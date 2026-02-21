package gcp

import (
	"context"
	"sync"

	"github.com/gourmet-guide/backend/internal/domain"
)

// SessionStore persists agent session metadata.
type SessionStore interface {
	SavePrompt(ctx context.Context, sessionID, prompt string) error
	SaveSession(ctx context.Context, session domain.ConciergeSession) error
	LoadSession(ctx context.Context, sessionID string) (domain.ConciergeSession, error)
	SaveMenuSafetyMetadata(ctx context.Context, restaurantID string, items []domain.MenuItem) error
	LoadMenuSafetyMetadata(ctx context.Context, restaurantID string) ([]domain.MenuItem, error)
	SaveImageReference(ctx context.Context, sessionID, imagePath string) error
	Close() error
}

// MemoryStore is local default storage for development and tests.
type MemoryStore struct {
	mu         sync.RWMutex
	sessions   map[string]domain.ConciergeSession
	menuByRest map[string][]domain.MenuItem
	images     map[string][]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		sessions:   map[string]domain.ConciergeSession{},
		menuByRest: map[string][]domain.MenuItem{},
		images:     map[string][]string{},
	}
}

func (m *MemoryStore) SavePrompt(_ context.Context, sessionID, prompt string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	session := m.sessions[sessionID]
	session.LastAssistantMsg = prompt
	m.sessions[sessionID] = session
	return nil
}

func (m *MemoryStore) SaveSession(_ context.Context, session domain.ConciergeSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[session.ID] = session
	return nil
}

func (m *MemoryStore) LoadSession(_ context.Context, sessionID string) (domain.ConciergeSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[sessionID], nil
}

func (m *MemoryStore) SaveMenuSafetyMetadata(_ context.Context, restaurantID string, items []domain.MenuItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.menuByRest[restaurantID] = append([]domain.MenuItem{}, items...)
	return nil
}

func (m *MemoryStore) LoadMenuSafetyMetadata(_ context.Context, restaurantID string) ([]domain.MenuItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := m.menuByRest[restaurantID]
	return append([]domain.MenuItem{}, items...), nil
}

func (m *MemoryStore) SaveImageReference(_ context.Context, sessionID, imagePath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.images[sessionID] = append(m.images[sessionID], imagePath)
	return nil
}

func (m *MemoryStore) Close() error { return nil }
