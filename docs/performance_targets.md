# Performance Targets and Release Gates

This document defines lightweight performance gates for the Hackathon MVP so we can catch regressions before demo deployments.

## Latency Budgets
- **Health check (`GET /healthz`)**: p95 <= 50ms in local test environment.
- **Session bootstrap (`POST /v1/sessions`)**: p95 <= 250ms with in-memory store.
- **Message turn (`POST /v1/sessions/:id/messages`)**: p95 <= 1200ms excluding model-network variance.
- **Realtime stream readiness (`GET /v1/sessions/:id/stream`)**: first `ready` event <= 200ms.

## Load Testing Gates
- Run a smoke load profile before promoting to staging/prod:
  - 20 virtual users for 2 minutes on `/healthz` and `/v1/sessions`.
  - Error rate must remain under 1%.
  - p95 for `/healthz` must remain <= 75ms.
- Any failed gate blocks release until regression is understood or threshold is consciously updated.

## CI / Local Enforcement
- Backend test suite includes `TestHealthEndpointLatencyBudget` as an automated latency gate.
- Stream resiliency and interruption behavior are covered by integration/unit tests in backend agent + HTTP handler suites.
- Frontend customer/admin journeys are validated with Vitest + Testing Library journey tests.
