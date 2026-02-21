package config

import "os"

// Config holds runtime settings for the API service.
type Config struct {
	Port            string
	ProjectID       string
	GeminiModel     string
	GoogleAPIKey    string
	FirestoreDBName string
	Region          string
}

// Load reads environment variables.
func Load() (Config, error) {
	cfg := Config{
		Port:            getenv("PORT", "8080"),
		ProjectID:       getenv("GOOGLE_CLOUD_PROJECT", "local-dev"),
		GeminiModel:     getenv("GEMINI_MODEL", "gemini-2.0-flash-live-001"),
		GoogleAPIKey:    os.Getenv("GOOGLE_API_KEY"),
		FirestoreDBName: getenv("FIRESTORE_DATABASE", "(default)"),
		Region:          getenv("GOOGLE_CLOUD_LOCATION", "us-central1"),
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
