# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]
### Added
- Added Cloud Storage provisioning in Terraform for menu image handling and exposed bucket output.
- Added Cloud Run cost controls in Terraform (min instance 0, high concurrency, lower memory target).
- Added backend runtime cost controls: relevant-menu-item limiting and in-memory prompt-response caching.

### Changed
- Refactored architecture/docs to the lean hackathon stack: Cloud Run + Firestore + Cloud Storage + Gemini on Vertex AI.
- Updated execution plan to remove Cloud SQL/Memorystore assumptions for MVP and align with cost-first delivery.
- Updated secrets guidance to prefer identity-based cloud auth and keep API keys local/optional.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-02-21
### Added
- Initialized repository for the GourmetGuide product roadmap.
- Added project guidance files: `AGENTS.md`, `execution_plan.md`, and `claude.md`.
- Added initial monorepo structure for:
  - `backend/` (Go services)
  - `frontend/` (React application)
  - `infra/` (Terraform infrastructure as code)

### Changed
- Updated `agents.md` to explicitly require checking `claude.md` before code generation.
