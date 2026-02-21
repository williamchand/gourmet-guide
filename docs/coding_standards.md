# GourmetGuide Coding Standards

## Cross-cutting
- Keep functions explicit and readable; avoid clever shortcuts.
- Eliminate duplicated logic quickly (DRY-first mindset).
- Treat allergen safety paths as high-risk: fail closed when uncertain.
- Every non-trivial change must include tests.

## Backend (Go)
- Use context-aware APIs (`context.Context`) for all external calls.
- Keep HTTP handlers thin; place domain rules in `internal/` packages.
- Return actionable error messages and wrap upstream failures.

## Frontend (React/JS)
- Keep UI state local unless shared across views.
- Separate display components from data-fetching concerns.
- Prefer semantic UI labels for accessibility.

## Infrastructure (Terraform)
- Keep modules/environment inputs explicit via `variables.tf`.
- Use managed services where possible before self-hosting.
- Keep IAM permissions least-privilege by default.
