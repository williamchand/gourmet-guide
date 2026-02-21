package agent

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gourmet-guide/backend/internal/domain"
	"github.com/gourmet-guide/backend/internal/gcp"
)

type blockingClient struct {
	started chan struct{}
}

func (b *blockingClient) Generate(ctx context.Context, _, _ string) (string, error) {
	close(b.started)
	<-ctx.Done()
	return "", ctx.Err()
}

func TestSendMessageReturnsInterruptionNoticeWhenRuntimeIsCanceled(t *testing.T) {
	t.Parallel()

	store := gcp.NewMemoryStore()
	runtime := NewRuntime("gemini", store)
	client := &blockingClient{started: make(chan struct{})}
	runtime.client = client
	service := NewConciergeService(store, gcp.NewMemoryImageStore(), runtime)

	_, err := service.SaveMenuItems(context.Background(), "rest-1", []domain.MenuItem{{Name: "Safe Bowl", Tags: []string{"vegan"}}})
	if err != nil {
		t.Fatalf("save menu: %v", err)
	}
	session, err := service.StartSession(context.Background(), "rest-1", nil, nil)
	if err != nil {
		t.Fatalf("start session: %v", err)
	}

	errCh := make(chan error, 1)
	replyCh := make(chan string, 1)
	go func() {
		reply, callErr := service.SendMessage(context.Background(), session.ID, "safe dinner?")
		if callErr != nil {
			errCh <- callErr
			return
		}
		replyCh <- reply
	}()

	select {
	case <-client.started:
	case <-time.After(2 * time.Second):
		t.Fatal("runtime client did not start")
	}

	if err := service.InterruptSession(context.Background(), session.ID); err != nil {
		t.Fatalf("interrupt session: %v", err)
	}

	select {
	case err := <-errCh:
		t.Fatalf("send message returned error: %v", err)
	case reply := <-replyCh:
		if reply != "response interrupted, ready for your next request" {
			t.Fatalf("unexpected interruption reply: %q", reply)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for interruption response")
	}

	updated, err := service.GetSession(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("load session: %v", err)
	}
	if updated.Status != domain.SessionStatusInterrupted {
		t.Fatalf("expected interrupted session status, got %s", updated.Status)
	}
}

func TestRuntimeRespondPropagatesCanceledContext(t *testing.T) {
	t.Parallel()
	store := gcp.NewMemoryStore()
	runtime := NewRuntime("gemini", store)
	runtime.client = &blockingClient{started: make(chan struct{})}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := runtime.Respond(ctx, "session-1", "hello", []string{"Soup"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled error, got %v", err)
	}
}
