package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/gourmet-guide/backend/internal/gcp"
)

const (
	maxMenuItemsDefault = 8
	maxItemLength       = 80
)

// Runtime uses Gemini-model-compatible clients and persists session activity.
type Runtime struct {
	modelName    string
	client       Client
	store        gcp.SessionStore
	maxMenuItems int

	mu    sync.RWMutex
	cache map[string]string
}

func NewRuntime(modelName string, store gcp.SessionStore) *Runtime {
	return &Runtime{
		modelName:    modelName,
		client:       newDefaultClient(),
		store:        store,
		maxMenuItems: maxMenuItemsDefault,
		cache:        map[string]string{},
	}
}

func (r *Runtime) Respond(ctx context.Context, sessionID, prompt string, menuItems []string) (string, error) {
	cleanPrompt, err := validatePrompt(prompt)
	if err != nil {
		return "", err
	}

	modelInput := r.buildModelInput(cleanPrompt, menuItems)
	if cachedReply, ok := r.cachedReply(modelInput); ok {
		if err := r.store.SavePrompt(ctx, sessionID, cleanPrompt); err != nil {
			return "", err
		}
		return cachedReply, nil
	}

	reply, err := r.client.Generate(ctx, r.modelName, modelInput)
	if err != nil {
		return "", err
	}

	r.cacheReply(modelInput, reply)
	if err := r.store.SavePrompt(ctx, sessionID, cleanPrompt); err != nil {
		return "", err
	}

	return reply, nil
}

func (r *Runtime) buildModelInput(prompt string, menuItems []string) string {
	relevantMenuItems := selectRelevantMenuItems(menuItems, r.maxMenuItems)
	if len(relevantMenuItems) == 0 {
		return prompt
	}

	return fmt.Sprintf(
		"%s\n\nOnly use these relevant menu options for reasoning:\n- %s",
		prompt,
		strings.Join(relevantMenuItems, "\n- "),
	)
}

func selectRelevantMenuItems(menuItems []string, maxItems int) []string {
	result := make([]string, 0, maxItems)
	for _, menuItem := range menuItems {
		if len(result) == maxItems {
			break
		}
		trimmed := strings.TrimSpace(menuItem)
		if trimmed == "" {
			continue
		}
		if len(trimmed) > maxItemLength {
			trimmed = trimmed[:maxItemLength]
		}
		result = append(result, trimmed)
	}
	return result
}

func (r *Runtime) cachedReply(modelInput string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	reply, ok := r.cache[modelInput]
	return reply, ok
}

func (r *Runtime) cacheReply(modelInput, reply string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[modelInput] = reply
}
