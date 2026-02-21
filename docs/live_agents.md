# Live Agents 🗣️ Blueprint

## Objective
Build a real-time, interruption-friendly GourmetGuide agent for audio + vision interactions on Google Cloud.

## Hackathon Compliance Requirements
All implementations in this repository must satisfy the following:
1. **Leverage a Gemini model**.
2. **Use Gemini Live API or ADK (Agent Development Kit)** for the agent runtime.
3. **Build with Google GenAI SDK or ADK**.
4. **Use at least one Google Cloud service** for hosting or data operations.

## Recommended Initial Build (Opinionated + Cost-Aware)
- **Runtime:** Gemini Live API with Google GenAI SDK through Vertex AI.
- **Hosting:** Cloud Run (min instances 0, concurrency 80, 512Mi memory target).
- **Session/Data:** Firestore for session memory and menu safety metadata.
- **Media:** Cloud Storage for menu images.

Avoid for hackathon MVP unless strictly required:
- GKE (always-on cluster spend)
- Cloud SQL (minimum instance baseline)
- Memorystore (always-on cost)

## Real-Time Interaction Scope
### Audio (Mandatory in MVP)
- Full-duplex voice streaming.
- User interruption support (cancel current model turn, restart on fresh audio input).
- Partial transcript updates and low-latency assistant responses.

### Vision (Mandatory in MVP)
- User can provide menu/dish image.
- Agent extracts candidate dish text.
- Agent cross-checks allergens against structured menu/allergen data.
- Agent returns safety explanation + safer alternatives when risk detected.


## Vision API Modes
- **Live user session mode (default):** text/audio interaction without mandatory image input.
- **Restaurant onboarding mode (alternative):** upload a menu image to auto-extract draft menu items and reduce manual restaurant data entry.
- Recommended endpoint for onboarding extraction: `POST /v1/restaurants/{restaurantId}/menu-extraction`.

## MVP Acceptance Criteria
- User can state allergies by voice and receive filtered recommendations.
- User can interrupt while the agent is responding, and the system switches context immediately.
- User can ask “Is this safe for me?” on a menu image and receive a reasoned answer.
- Service is deployed on Google Cloud and uses at least one managed GCP service.

## Gemini Cost Controls (Required)
- Limit output tokens for every model call.
- Cache repeat-safe filtering answers.
- Send only relevant menu items; never send entire menu/database per turn.
- Prefer structured tool outputs to long free-form reasoning where possible.

## Reference Lean Architecture
Frontend (static hosting) → Cloud Run backend → Gemini Live (Vertex AI) → Firestore + Cloud Storage
