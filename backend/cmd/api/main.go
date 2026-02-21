package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gourmet-guide/backend/internal/agent"
	"github.com/gourmet-guide/backend/internal/config"
	"github.com/gourmet-guide/backend/internal/gcp"
)

type requestBody struct {
	SessionID string   `json:"sessionId"`
	Prompt    string   `json:"prompt"`
	MenuItems []string `json:"menuItems"`
}

type responseBody struct {
	Reply string `json:"reply"`
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	store := gcp.NewMemoryStore()
	defer store.Close()

	runtime := agent.NewRuntime(cfg.GeminiModel, store)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/v1/agent/respond", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req requestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()

		reply, err := runtime.Respond(ctx, req.SessionID, req.Prompt, req.MenuItems)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(responseBody{Reply: reply})
	})

	log.Printf("backend listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
