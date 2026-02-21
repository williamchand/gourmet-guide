package main

import (
	"log"
	"net/http"

	"github.com/gourmet-guide/backend/internal/agent"
	"github.com/gourmet-guide/backend/internal/config"
	"github.com/gourmet-guide/backend/internal/gcp"
	httphandler "github.com/gourmet-guide/backend/internal/handler/http"
	"github.com/gourmet-guide/backend/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	store := gcp.NewMemoryStore()
	defer store.Close()

	runtime := agent.NewRuntime(cfg.GeminiModel, store)
	concierge := agent.NewConciergeService(store, gcp.NewMemoryImageStore(), runtime)
	app := service.NewConciergeApp(concierge)
	handler := httphandler.NewHandler(app)

	log.Printf("backend listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler.Routes()); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
