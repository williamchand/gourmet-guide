package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gourmet-guide/backend/internal/agent"
	"github.com/gourmet-guide/backend/internal/gcp"
	"github.com/gourmet-guide/backend/internal/service"
)

func testServer() http.Handler {
	store := gcp.NewMemoryStore()
	runtime := agent.NewRuntime("gemini", store)
	concierge := agent.NewConciergeService(store, gcp.NewMemoryImageStore(), runtime)
	app := service.NewConciergeApp(concierge)
	return NewHandler(app).Routes()
}

func TestSessionAndMenuExtractionRoutes(t *testing.T) {
	t.Parallel()
	router := testServer()

	startPayload := map[string]any{
		"restaurantId":   "rest-e2e",
		"hardAllergens":  []string{"peanut"},
		"preferenceTags": []string{"vegan"},
		"menuItems": []map[string]any{
			{"name": "Tofu Bowl", "tags": []string{"vegan"}},
		},
	}
	body, _ := json.Marshal(startPayload)
	req := httptest.NewRequest(http.MethodPost, "/v1/sessions", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from start session, got %d (%s)", rec.Code, rec.Body.String())
	}

	extractionPayload := map[string]string{
		"fileName": "menu.txt",
		"base64":   base64.StdEncoding.EncodeToString([]byte("Tofu Bowl\nPeanut Curry")),
	}
	body, _ = json.Marshal(extractionPayload)
	req = httptest.NewRequest(http.MethodPost, "/v1/restaurants/rest-e2e/menu-extraction", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from menu extraction, got %d (%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Vision extraction is optional") {
		t.Fatalf("expected extraction guidance note, got %s", rec.Body.String())
	}
}

func TestStreamEndpointSendsReadyEventAndHandlesClientCancel(t *testing.T) {
	t.Parallel()
	router := testServer()

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/v1/sessions/unknown/stream", nil).WithContext(ctx)
	rec := httptest.NewRecorder()
	done := make(chan struct{})

	go func() {
		router.ServeHTTP(rec, req)
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("stream handler did not exit after client cancellation")
	}

	if !strings.Contains(rec.Body.String(), "event: ready") {
		t.Fatalf("expected ready event in stream response, got %s", rec.Body.String())
	}
}
