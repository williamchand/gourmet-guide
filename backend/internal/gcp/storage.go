package gcp

import (
	"context"
	"fmt"
	"path"
	"sync"
)

// ImageStore stores uploaded menu images for vision safety checks.
type ImageStore interface {
	SaveSessionImage(ctx context.Context, sessionID, fileName string, content []byte) (string, error)
}

type MemoryImageStore struct {
	mu      sync.Mutex
	objects map[string][]byte
}

func NewMemoryImageStore() *MemoryImageStore {
	return &MemoryImageStore{objects: map[string][]byte{}}
}

func (s *MemoryImageStore) SaveSessionImage(_ context.Context, sessionID, fileName string, content []byte) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := path.Join(sessionID, fileName)
	s.objects[key] = append([]byte{}, content...)
	return fmt.Sprintf("memory://%s", key), nil
}
