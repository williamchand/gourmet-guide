# GourmetGuide

A Friendly, Allergen-Aware AI Restaurant Concierge with Smart Combo Recommendations.

## Repository Initialization
This repository is initialized as a monorepo with the following stack:
- **Backend:** Golang
- **Frontend:** React
- **Infrastructure as Code:** Terraform

## Live Agents Focus
This project includes a dedicated **Live Agents 🗣️** track for real-time interaction (audio + vision):
- Agent runtime must use **Gemini Live API or ADK**
- Implementation must use **Google GenAI SDK or ADK**
- Deployment must run on **Google Cloud** and use at least one managed GCP service

## Structure
- `backend/` — Go services and business logic
- `frontend/` — React applications (customer and admin experiences)
- `infra/` — Terraform modules and environment configurations
- `docs/` — architecture and planning docs

## Planning Documents
- `changelog.md` — release history
- `execution_plan.md` — staged implementation roadmap
- `agents.md` — project guardrails and contributor instructions
- `claude.md` — structured review prompt for collaborative planning
- `docs/live_agents.md` — live agent scope, constraints, and MVP acceptance criteria
