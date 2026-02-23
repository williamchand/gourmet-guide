package http

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
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

func createSession(t *testing.T, router http.Handler) string {
	t.Helper()
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
	var parsed struct {
		Session struct {
			ID string `json:"id"`
		} `json:"session"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("parse start session response: %v", err)
	}
	if parsed.Session.ID == "" {
		t.Fatal("expected non-empty session id")
	}
	return parsed.Session.ID
}

func TestSessionAndMenuExtractionRoutes(t *testing.T) {
	t.Parallel()
	router := testServer()
	_ = createSession(t, router)

	extractionPayload := map[string]string{
		"fileName": "menu.txt",
		"base64":   base64.StdEncoding.EncodeToString([]byte("Tofu Bowl\nPeanut Curry")),
	}
	body, _ := json.Marshal(extractionPayload)
	req := httptest.NewRequest(http.MethodPost, "/v1/restaurants/rest-e2e/menu-extraction", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from menu extraction, got %d (%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Vision extraction is optional") {
		t.Fatalf("expected extraction guidance note, got %s", rec.Body.String())
	}
}

func TestVoiceStreamingConfigEndpoint(t *testing.T) {
	t.Parallel()
	router := testServer()

	req := httptest.NewRequest(http.MethodGet, "/v1/realtime/voice-config", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from voice config, got %d (%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "gemini-2.5-flash-native-audio-preview-12-2025") {
		t.Fatalf("expected gemini native audio model, got %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "\"response_modalities\":[\"AUDIO\"]") {
		t.Fatalf("expected audio response modality, got %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "\"speech_config\":") {
		t.Fatalf("expected speech_config in voice config, got %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "\"input_mime_type\":\"audio/pcm\"") {
		t.Fatalf("expected audio/pcm input mime type, got %s", rec.Body.String())
	}
}

func TestRealtimeWebSocketFlow(t *testing.T) {
	t.Parallel()
	router := testServer()
	sessionID := createSession(t, router)

	srv := httptest.NewServer(router)
	defer srv.Close()

	host := strings.TrimPrefix(srv.URL, "http://")
	conn, rw := dialWS(t, host, "/ws/user-1/"+sessionID)
	defer conn.Close()

	ready := readWSText(t, rw)
	if !strings.Contains(ready, "\"type\":\"ready\"") {
		t.Fatalf("expected ready event, got %s", ready)
	}

	writeWSText(t, rw, `{"type":"close"}`)
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

func dialWS(t *testing.T, host, path string) (net.Conn, *bufio.ReadWriter) {
	t.Helper()
	conn, err := net.Dial("tcp", host)
	if err != nil {
		t.Fatalf("dial tcp: %v", err)
	}
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	keyBytes := make([]byte, 16)
	if _, err := rand.Read(keyBytes); err != nil {
		t.Fatalf("random key: %v", err)
	}
	key := base64.StdEncoding.EncodeToString(keyBytes)

	req := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: %s\r\nSec-WebSocket-Version: 13\r\n\r\n", path, host, key)
	if _, err := rw.WriteString(req); err != nil {
		t.Fatalf("write handshake: %v", err)
	}
	if err := rw.Flush(); err != nil {
		t.Fatalf("flush handshake: %v", err)
	}

	status, err := rw.ReadString('\n')
	if err != nil {
		t.Fatalf("read status: %v", err)
	}
	if !strings.Contains(status, "101") {
		t.Fatalf("expected 101 status, got %s", status)
	}
	for {
		line, err := rw.ReadString('\n')
		if err != nil {
			t.Fatalf("read header: %v", err)
		}
		if line == "\r\n" {
			break
		}
	}

	return conn, rw
}

func writeWSText(t *testing.T, rw *bufio.ReadWriter, payload string) {
	t.Helper()
	data := []byte(payload)
	header := []byte{0x81, byte(0x80 | len(data))}
	mask := []byte{0x11, 0x22, 0x33, 0x44}
	masked := make([]byte, len(data))
	for i := range data {
		masked[i] = data[i] ^ mask[i%4]
	}
	if _, err := rw.Write(header); err != nil {
		t.Fatalf("write ws header: %v", err)
	}
	if _, err := rw.Write(mask); err != nil {
		t.Fatalf("write ws mask: %v", err)
	}
	if _, err := rw.Write(masked); err != nil {
		t.Fatalf("write ws payload: %v", err)
	}
	if err := rw.Flush(); err != nil {
		t.Fatalf("flush ws payload: %v", err)
	}
}

func readWSText(t *testing.T, rw *bufio.ReadWriter) string {
	t.Helper()
	head := make([]byte, 2)
	if _, err := io.ReadFull(rw, head); err != nil {
		t.Fatalf("read ws header: %v", err)
	}
	length := int(head[1] & 0x7F)
	payload := make([]byte, length)
	if _, err := io.ReadFull(rw, payload); err != nil {
		t.Fatalf("read ws payload: %v", err)
	}
	return string(payload)
}
