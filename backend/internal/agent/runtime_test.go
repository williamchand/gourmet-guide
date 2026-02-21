package agent

import (
	"context"
	"testing"

	"github.com/gourmet-guide/backend/internal/gcp"
)

type fakeClient struct {
	calls int
}

func (f *fakeClient) Generate(_ context.Context, _, prompt string) (string, error) {
	f.calls++
	return "ok: " + prompt, nil
}

func TestPromptValidation(t *testing.T) {
	t.Parallel()

	_, err := validatePrompt("   ")
	if err == nil {
		t.Fatal("expected error for empty prompt")
	}

	trimmed, err := validatePrompt("  hello ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if trimmed != "hello" {
		t.Fatalf("expected trimmed prompt, got %q", trimmed)
	}
}

func TestSelectRelevantMenuItemsLimitsItems(t *testing.T) {
	t.Parallel()
	items := []string{" one ", "", "two", "three", "four"}
	selected := selectRelevantMenuItems(items, 2)
	if len(selected) != 2 {
		t.Fatalf("expected 2 items, got %d", len(selected))
	}
	if selected[0] != "one" || selected[1] != "two" {
		t.Fatalf("unexpected selection: %#v", selected)
	}
}

func TestRespondUsesCacheToReduceModelCalls(t *testing.T) {
	t.Parallel()
	store := gcp.NewMemoryStore()
	runtime := NewRuntime("gemini-2.0-flash-live-001", store)
	fake := &fakeClient{}
	runtime.client = fake

	menu := []string{"A", "B", "C"}
	_, err := runtime.Respond(context.Background(), "session-1", "safe options?", menu)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = runtime.Respond(context.Background(), "session-1", "safe options?", menu)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.calls != 1 {
		t.Fatalf("expected one model call due to cache, got %d", fake.calls)
	}
}
