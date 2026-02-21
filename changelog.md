# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]
### Added
- Added `docs/live_agents.md` to define Live Agents (audio/vision) scope, mandatory technology constraints, and MVP acceptance criteria.

### Changed
- Updated `README.md` to include Live Agents requirements and reference to `docs/live_agents.md`.
- Updated `execution_plan.md` with a Live Agents compliance section and vision/interruption-specific milestones.

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
