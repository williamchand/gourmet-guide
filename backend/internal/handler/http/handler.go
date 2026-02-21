package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gourmet-guide/backend/internal/domain"
	"github.com/gourmet-guide/backend/internal/service"
)

type Handler struct {
	app *service.ConciergeApp
}

func NewHandler(app *service.ConciergeApp) *Handler {
	return &Handler{app: app}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.handleHealth)
	mux.HandleFunc("/v1/sessions", h.handleSessions)
	mux.HandleFunc("/v1/sessions/", h.handleSessionByID)
	mux.HandleFunc("/v1/restaurants/", h.handleRestaurantRoutes)
	return mux
}

type startSessionRequest struct {
	RestaurantID   string            `json:"restaurantId"`
	HardAllergens  []domain.Allergen `json:"hardAllergens"`
	PreferenceTags []string          `json:"preferenceTags"`
	MenuItems      []domain.MenuItem `json:"menuItems"`
}

type sendMessageRequest struct {
	Prompt string `json:"prompt"`
}

type imageUploadRequest struct {
	FileName string `json:"fileName"`
	Base64   string `json:"base64"`
}

type menuTaggingRequest struct {
	MenuItems []domain.MenuItem `json:"menuItems"`
}

type sessionStartResponse struct {
	Session           domain.ConciergeSession `json:"session"`
	SuggestedMenuTags []domain.MenuItem       `json:"suggestedMenuItems"`
}

type menuExtractionResponse struct {
	ImagePath string            `json:"imagePath"`
	MenuItems []domain.MenuItem `json:"menuItems"`
	Note      string            `json:"note"`
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *Handler) handleSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req startSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := h.app.StartSession(ctx, service.StartSessionInput{
		RestaurantID:   req.RestaurantID,
		HardAllergens:  req.HardAllergens,
		PreferenceTags: req.PreferenceTags,
		MenuItems:      req.MenuItems,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, sessionStartResponse{Session: result.Session, SuggestedMenuTags: result.SuggestedMenuItems})
}

func (h *Handler) handleSessionByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v1/sessions/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	sessionID := parts[0]

	if len(parts) == 1 && r.Method == http.MethodGet {
		session, err := h.app.GetSession(r.Context(), sessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, session)
		return
	}
	if len(parts) == 1 && r.Method == http.MethodDelete {
		if err := h.app.EndSession(r.Context(), sessionID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if len(parts) == 2 && parts[1] == "messages" && r.Method == http.MethodPost {
		var req sendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		reply, err := h.app.SendMessage(r.Context(), sessionID, req.Prompt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, map[string]string{"reply": reply})
		return
	}
	if len(parts) == 2 && parts[1] == "interrupt" && r.Method == http.MethodPost {
		if err := h.app.InterruptSession(r.Context(), sessionID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		return
	}
	if len(parts) == 2 && parts[1] == "ws" && r.Method == http.MethodGet {
		http.Error(w, "websocket transport endpoint reserved for Gemini Live clients", http.StatusUpgradeRequired)
		return
	}
	if len(parts) == 2 && parts[1] == "stream" && r.Method == http.MethodGet {
		handleRealtimeStream(w, r, h.app, sessionID)
		return
	}
	http.NotFound(w, r)
}

func (h *Handler) handleRestaurantRoutes(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v1/restaurants/"), "/")
	if len(parts) != 2 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	restaurantID := parts[0]
	if parts[1] == "menu-tags" && r.Method == http.MethodPost {
		var req menuTaggingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		enriched, err := h.app.TagMenuItems(r.Context(), restaurantID, req.MenuItems)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]any{"menuItems": enriched, "note": "Tags were auto-suggested to simplify allergy/diet filters for business owners."})
		return
	}
	if parts[1] != "menu-extraction" || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var req imageUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	content, err := base64.StdEncoding.DecodeString(req.Base64)
	if err != nil {
		http.Error(w, "invalid base64 image", http.StatusBadRequest)
		return
	}
	result, err := h.app.ExtractMenuFromImage(r.Context(), restaurantID, req.FileName, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, menuExtractionResponse{
		ImagePath: result.ImagePath,
		MenuItems: result.MenuItems,
		Note:      "Vision extraction is optional for onboarding; for live interaction, use text/audio session APIs.",
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func handleRealtimeStream(w http.ResponseWriter, r *http.Request, app *service.ConciergeApp, sessionID string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("event: ready\ndata: stream-open\n\n"))
	flusher.Flush()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			session, err := app.GetSession(r.Context(), sessionID)
			if err != nil {
				return
			}
			payload, _ := json.Marshal(session)
			_, _ = w.Write([]byte("event: session\ndata: " + string(payload) + "\n\n"))
			flusher.Flush()
		}
	}
}
