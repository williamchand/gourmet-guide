package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gourmet-guide/backend/internal/domain"
	"github.com/gourmet-guide/backend/internal/service"
)

type Handler struct {
	app *service.ConciergeApp
}

type voiceStreamingConfigResponse struct {
	Model  string `json:"model"`
	Config struct {
		ResponseModalities           []string       `json:"response_modalities"`
		SystemInstruction            string         `json:"system_instruction"`
		SpeechConfig                 map[string]any `json:"speech_config"`
		ManualActivitySignalsEnabled bool           `json:"manual_activity_signals_enabled"`
		EnableProactiveAudio         bool           `json:"enable_proactive_audio"`
		EnableAffectiveDialog        bool           `json:"enable_affective_dialog"`
		SessionResumptionEnabled     bool           `json:"session_resumption_enabled"`
		ContextWindowCompressionOn   bool           `json:"context_window_compression_enabled"`
		ContextWindowCompression     map[string]int `json:"context_window_compression"`
		SaveLiveBlob                 bool           `json:"save_live_blob"`
		SupportCFC                   bool           `json:"support_cfc"`
		MaxLLMCalls                  *int           `json:"max_llm_calls"`
		CustomMetadata               map[string]any `json:"custom_metadata"`
	} `json:"config"`
	Audio struct {
		Format            string `json:"format"`
		Channels          int    `json:"channels"`
		SendSampleRate    int    `json:"send_sample_rate"`
		ReceiveSampleRate int    `json:"receive_sample_rate"`
		ChunkSize         int    `json:"chunk_size"`
		InputMimeType     string `json:"input_mime_type"`
		OutputMimeType    string `json:"output_mime_type"`
	} `json:"audio"`
}

func NewHandler(app *service.ConciergeApp) *Handler {
	return &Handler{app: app}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.handleHealth)
	mux.HandleFunc("/v1/realtime/voice-config", h.handleVoiceStreamingConfig)
	mux.HandleFunc("/ws/", h.handleRealtimeWebSocket)
	mux.HandleFunc("/v1/sessions", h.handleSessions)
	mux.HandleFunc("/v1/sessions/", h.handleSessionByID)
	mux.HandleFunc("/v1/restaurants/", h.handleRestaurantRoutes)
	return mux
}

func (h *Handler) handleVoiceStreamingConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	response := voiceStreamingConfigResponse{Model: getenv("GEMINI_MODEL", "gemini-2.5-flash-native-audio-preview-12-2025")}
	response.Config.ResponseModalities = []string{"AUDIO"}
	response.Config.SystemInstruction = "You are a helpful and friendly AI assistant."
	response.Config.SpeechConfig = map[string]any{
		"voice_name":    getenv("VOICE_NAME", "Aoede"),
		"language_code": getenv("VOICE_LANGUAGE_CODE", "en-US"),
	}
	response.Config.ManualActivitySignalsEnabled = getenvBool("ENABLE_MANUAL_ACTIVITY_SIGNALS", false)
	response.Config.EnableProactiveAudio = getenvBool("ENABLE_PROACTIVE_AUDIO", false)
	response.Config.EnableAffectiveDialog = getenvBool("ENABLE_AFFECTIVE_DIALOG", false)
	response.Config.SessionResumptionEnabled = getenvBool("ENABLE_SESSION_RESUMPTION", true)
	response.Config.ContextWindowCompressionOn = getenvBool("ENABLE_CONTEXT_WINDOW_COMPRESSION", false)
	response.Config.ContextWindowCompression = map[string]int{
		"trigger_tokens": getenvInt("CONTEXT_COMPRESSION_TRIGGER_TOKENS", 100000),
		"target_tokens":  getenvInt("CONTEXT_COMPRESSION_TARGET_TOKENS", 80000),
	}
	response.Config.SaveLiveBlob = getenvBool("SAVE_LIVE_BLOB", false)
	response.Config.SupportCFC = getenvBool("SUPPORT_CFC", false)
	response.Config.MaxLLMCalls = getenvOptionalInt("MAX_LLM_CALLS")
	response.Config.CustomMetadata = map[string]any{
		"app":               "gourmet-guide-bidi",
		"transport":         "websocket",
		"response_modality": "AUDIO",
	}
	response.Audio.Format = "pcm16"
	response.Audio.Channels = 1
	response.Audio.SendSampleRate = 16000
	response.Audio.ReceiveSampleRate = 24000
	response.Audio.ChunkSize = 1024
	response.Audio.InputMimeType = "audio/pcm"
	response.Audio.OutputMimeType = "audio/pcm"

	writeJSON(w, response)
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
		h.handleRealtimeWebSocketSession(w, r, sessionID)
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

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func getenvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if value == "" {
		return fallback
	}
	return value == "1" || value == "true" || value == "yes" || value == "on"
}

func getenvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getenvOptionalInt(key string) *int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &parsed
}
