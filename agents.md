# AGENTS.md

## Project Scope
This file applies to the full repository.

## Product Context
GourmetGuide is an allergen-aware AI restaurant concierge that combines:
- Go backend services
- React frontend experiences
- Terraform-managed Google Cloud infrastructure

## Engineering Guardrails
- Prefer explicit, readable code over clever abstractions.
- Be aggressive about DRY violations and remove duplicated logic early.
- Add tests for all non-trivial behavior (unit first, integration where needed).
- Handle edge cases deliberately, especially for allergen and safety workflows.
- Keep architecture "engineered enough": avoid both fragile shortcuts and premature complexity.

## Repository Conventions
- `backend/`: Go source code (API, tool orchestration, and business logic)
- `frontend/`: React app (customer + admin experiences)
- `infra/`: Terraform modules and environment stacks
- `docs/`: architecture and operational documentation (create as needed)

## Delivery Expectations
When proposing design or code changes:
1. Provide concrete tradeoffs.
2. Give an opinionated recommendation.
3. Ask for user confirmation before major direction changes.

## Code Generation Instruction
Before generating code or proposing implementation details, review `claude.md` and follow its interaction pattern and review constraints.

## Testing Expectations
- Backend: `go test ./...`
- Frontend: `npm test` (or configured test runner)
- Infra: `terraform fmt -check`, `terraform validate`
