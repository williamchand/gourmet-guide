# GourmetGuide

A friendly, allergen-aware AI restaurant concierge with smart combo recommendations.

## Live Agents Compliance (Execution Plan 0-1)
This repository includes a cost-aware foundation for the first delivery phase:
- ✅ Gemini model usage for all agentic runtime calls (`gemini-2.0-flash-live-001` default).
- ✅ Agent runtime scaffold using Google GenAI SDK, with Vertex AI path for Gemini in cloud builds.
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
```

### Frontend
```bash
cd frontend
npm install
npm run lint
npm test
npm run build
```

### Infrastructure
```bash
cd infra
terraform init -backend=false
terraform fmt -check
terraform validate
```

## Governance and Standards
- Coding standards: `docs/coding_standards.md`
- Secrets strategy: `docs/secrets.md`
- PR checklist template: `.github/pull_request_template.md`
- CI workflow: `.github/workflows/ci.yml`
