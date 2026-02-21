# GourmetGuide Execution Plan

## 0. Live Agents Hackathon Compliance
- [ ] Use a Gemini model for all agentic interactions.
- [ ] Implement agent runtime with **Gemini Live API or ADK**.
- [ ] Use **Google GenAI SDK or ADK** in implementation.
- [ ] Deploy on Google Cloud and integrate at least one GCP managed service.

## 1. Foundation (Repository + Tooling)
- [ ] Establish monorepo conventions (`backend/`, `frontend/`, `infra/`).
- [ ] Add linting, formatting, and CI workflows.
- [ ] Define coding standards and PR template.
- [ ] Configure secret management approach for local + cloud environments.

## 2. Backend (Go)
### 2.1 Core API
- [ ] Build HTTP/WebSocket API surface for realtime concierge sessions.
- [ ] Implement session lifecycle management and interruption handling.
- [ ] Define domain models for restaurants, menu items, allergens, and combos.

### 2.2 Data + Integrations
- [ ] Integrate Cloud SQL (Postgres) for menu and combo data.
- [ ] Integrate Memorystore (Redis) for conversational session state.
- [ ] Add Gemini Live + tool-calling orchestration.

### 2.3 Safety
- [ ] Implement hard allergen filters and preference ranking.
- [ ] Add cross-contamination policy checks.
- [ ] Add fallback handling and disclaimer responses for high-risk queries.

## 3. Frontend (React)
### 3.1 Customer Experience
- [ ] Build microphone-driven realtime assistant UI.
- [ ] Add transcript stream and recommendation panel.
- [ ] Support allergy profile preferences and updates mid-session.
- [ ] Add menu image upload/capture flow for vision safety checks.

### 3.2 Admin Experience
- [ ] Create menu management dashboard.
- [ ] Add allergen/ingredient tagging forms.
- [ ] Add combo builder and preview tools.

## 4. Infrastructure (Terraform on Google Cloud)
- [ ] Provision Cloud Run for backend services.
- [ ] Provision Cloud SQL and Memorystore.
- [ ] Provision Cloud Storage for menu/dish images.
- [ ] Configure IAM, networking, and least-privilege service accounts.
- [ ] Add environment stacks (`dev`, `staging`, `prod`).

## 5. Quality and Reliability
- [ ] Unit tests across backend domain logic and filtering.
- [ ] Integration tests for database/tool integration paths.
- [ ] E2E tests for key customer/admin journeys.
- [ ] Add interruption-path and streaming resiliency tests.
- [ ] Performance targets: latency budgets and load-testing gates.

## 6. Release Plan
### Phase 1 (Hackathon MVP)
- [ ] Voice interaction + allergy filtering
- [ ] Vision-based safety checks
- [ ] Basic combo recommendations
- [ ] Cloud deployment demo

### Phase 2
- [ ] Nutritional scoring
- [ ] Predictive upsell models
- [ ] Cross-contamination AI assistance

### Phase 3
- [ ] POS integration
- [ ] Loyalty personalization
- [ ] Advanced analytics dashboard
