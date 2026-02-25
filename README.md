# GourmetGuide

A friendly, allergen-aware AI restaurant concierge with smart combo recommendations.

## Live Agents Compliance (Execution Plan 0-1)
This repository includes a cost-aware foundation for the first delivery phase:
- ✅ Gemini model usage for all agentic runtime calls (`gemini-2.5-flash-native-audio-preview-12-2025` default).
- ✅ Agent runtime scaffold using Google GenAI SDK, with Vertex AI path for Gemini in cloud builds.
- ✅ Voice streaming contract endpoint (`GET /v1/realtime/voice-config`) aligned with Gemini Live audio settings (16kHz in / 24kHz out PCM).
- ✅ Cheap managed GCP baseline: Cloud Run + Firestore + Cloud Storage.
- ✅ Monorepo conventions and CI/tooling baseline across backend, frontend, and infra.

## Lean Hackathon Architecture
Frontend (static hosting) → Cloud Run backend → Gemini Live (Vertex AI) → Firestore + Cloud Storage

### Cost controls baked into this repo
- Cloud Run is configured for min instances `0`, concurrency `80`, and `512Mi` memory target.
- Backend runtime limits model context to relevant menu items only and caches repeated prompts.
- Firestore + Cloud Storage are used instead of always-on datastores for MVP.

## Monorepo Structure
- `backend/` — Go API and agent runtime integration (Google GenAI SDK + Firestore)
- `frontend/` — JS frontend workspace scaffold with lint/test/format scripts
- `infra/` — Terraform definitions for Google Cloud deployment
- `docs/` — coding standards, architecture context, and secrets strategy

## Quick Start
### Backend
```bash
cd backend
go test ./...
go run ./cmd/api
# generate ~100 Gemini-powered menu items + corresponding food images (requires GOOGLE_API_KEY)
GOOGLE_API_KEY=your_key GOOGLE_CLOUD_PROJECT=your_project GCS_BUCKET=your_bucket \
  go run -tags gcp ./cmd/seeddata --count 100 --out seed/output
# optional: keep local-only output without cloud writes
GOOGLE_API_KEY=your_key go run ./cmd/seeddata --count 100 --skip-gcs-upload --skip-firestore-write
```

### Frontend
```bash
cd frontend
npm install
npm run lint
npm test
npm run build
```

### Frontend deploy (GCS static hosting)
```bash
scripts/deploy_frontend_gcs.sh --bucket <your-bucket-name> --project <your-project-id>
```
Detailed guide: `docs/frontend_gcs_deploy.md`.


### Realtime voice websocket (Go backend)
The main Go backend now exposes realtime voice endpoints directly:
- `GET /v1/realtime/voice-config`
- `WS /ws/{user_id}/{session_id}`
- `WS /v1/sessions/{session_id}/ws`

Upstream websocket accepts binary audio frames (`audio/pcm;rate=16000`) or JSON messages (`text`, `audio`, `image`, `activity_start`, `activity_end`, `close`).
Manual activity signals require `ENABLE_MANUAL_ACTIVITY_SIGNALS=true`.
Frontend dev server proxies `/v1` and `/ws` to `http://localhost:8080` so the customer UI uses the Go backend realtime endpoints during local development.

### Infrastructure
```bash
cd infra
terraform init -backend=false
terraform fmt -check
terraform validate
```


### Automated deploy with GitHub Actions (Terraform + container registry)
Tag-based releases are fully automated using `.github/workflows/deploy.yml`:
- Builds and pushes backend Docker image to a container registry (Docker Hub, GitHub Packages, etc.) using the release tag as image tag (fallback: commit SHA).
- Runs `terraform apply` for backend infrastructure on Google Cloud.
- Builds frontend and syncs `frontend/dist` to a Google Cloud Storage bucket.

Trigger this workflow by creating a tag such as `v1.0.0` and pushing it to GitHub.

Set these repository **Variables** (all except the registry prefix are still GCP‑specific):
- `GCP_PROJECT_ID`
- `GCP_REGION`
- `DEPLOY_ENVIRONMENT` (for example `prod`)
- `FRONTEND_GCS_BUCKET`
- `DOCKER_REGISTRY` – the full image name to push, for example `docker.io/myorg/gourmet-guide-backend` or `ghcr.io/myorg/gourmet-guide-backend`.

Set these repository **Secrets**:
- `DOCKER_USERNAME` (for the registry)
- `DOCKER_PASSWORD` (or token)
- `GCP_WORKLOAD_IDENTITY_PROVIDER`
- `GCP_SERVICE_ACCOUNT`
> The service account behind `GCP_SERVICE_ACCOUNT` must allow Terraform resource management and GCS static-asset deployment.

## Local Emulator Workflow
Use `docker-compose.dev.yml` to run Firestore + Cloud Storage emulators for local seed/publish testing.

```bash
docker compose -f docker-compose.dev.yml up -d
```

Detailed steps and env vars: `docs/local_emulators.md`.

## Governance and Standards
- Coding standards: `docs/coding_standards.md`
- Secrets strategy: `docs/secrets.md`
- PR checklist template: `.github/pull_request_template.md`
- CI workflow: `.github/workflows/ci.yml`
