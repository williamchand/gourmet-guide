# GourmetGuide Execution Plan

## 0. Live Agents Hackathon Compliance
- [x] Use a Gemini model for all agentic interactions.
- [x] Implement agent runtime with **Gemini Live API or ADK**.
- [x] Use **Google GenAI SDK or ADK** in implementation.
- [x] Deploy on Google Cloud and integrate at least one GCP managed service.

## 1. Foundation (Repository + Tooling)
- [x] Establish monorepo conventions (`backend/`, `frontend/`, `infra/`).
- [x] Add linting, formatting, and CI workflows.
- [x] Define coding standards and PR template.
- [x] Configure secret management approach for local + cloud environments.

## 2. Backend (Go)
### 2.1 Core API
- [x] Build HTTP/WebSocket API surface for realtime concierge sessions.
- [x] Implement session lifecycle management and interruption handling.
- [x] Define domain models for restaurants, menu items, allergens, and combos.

### 2.2 Data + Integrations
- [x] Integrate Firestore for conversational session state and lightweight menu safety metadata.
- [x] Add Cloud Storage image workflows for vision safety checks.
- [x] Add Gemini Live + tool-calling orchestration.

### 2.3 Safety
- [x] Implement hard allergen filters and preference ranking.
- [x] Add cross-contamination policy checks.
- [x] Add fallback handling and disclaimer responses for high-risk queries.

## 3. Frontend (React)
### 3.1 Customer Experience
- [x] Build microphone-driven realtime assistant UI.
- [x] Add transcript stream and recommendation panel.
- [x] Support allergy profile preferences and updates mid-session.
- [x] Add menu image upload/capture flow for vision safety checks.

### 3.2 Admin Experience
- [x] Create menu management dashboard.
- [x] Add allergen/ingredient tagging forms.
- [x] Add combo builder and preview tools.

## 4. Infrastructure (Terraform on Google Cloud)
- [x] Provision Cloud Run for backend services.
- [x] Provision Firestore.
- [x] Provision Cloud Storage for menu/dish images.
- [x] Configure IAM, networking, and least-privilege service accounts.
- [x] Add environment stacks (`dev`, `staging`, `prod`).

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
