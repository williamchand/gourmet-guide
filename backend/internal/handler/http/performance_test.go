package http

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

func TestHealthEndpointLatencyBudget(t *testing.T) {
	t.Parallel()
	router := testServer()

	const samples = 40
	latencies := make([]time.Duration, 0, samples)
	for i := 0; i < samples; i++ {
		start := time.Now()
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("healthz request failed with status %d", rec.Code)
		}
		latencies = append(latencies, time.Since(start))
	}

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	p95 := latencies[int(float64(samples)*0.95)-1]
	if p95 > 50*time.Millisecond {
		t.Fatalf("p95 latency budget exceeded: got %s, max 50ms", p95)
	}
}
