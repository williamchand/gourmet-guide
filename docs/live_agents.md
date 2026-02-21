# Live Agents 🗣️ Blueprint

## Objective
Build a real-time, interruption-friendly GourmetGuide agent for audio + vision interactions on Google Cloud.

## Hackathon Compliance Requirements
All implementations in this repository must satisfy the following:
1. **Leverage a Gemini model**.
2. **Use Gemini Live API or ADK (Agent Development Kit)** for the agent runtime.
3. **Build with Google GenAI SDK or ADK**.
4. **Use at least one Google Cloud service** for hosting or data operations.

## Recommended Initial Build (Opinionated)
- **Runtime:** Gemini Live API with Google GenAI SDK.
- **Hosting:** Cloud Run.
- **State:** Memorystore (Redis) for live session context + interruption state.
- **Data:** Cloud SQL (Postgres) for menus/allergens.
- **Media:** Cloud Storage for menu images.

This path minimizes custom orchestration overhead while preserving a production-friendly path to ADK later if multi-agent composition becomes necessary.

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

## MVP Acceptance Criteria
- User can state allergies by voice and receive filtered recommendations.
- User can interrupt while the agent is responding, and the system switches context immediately.
- User can ask “Is this safe for me?” on a menu image and receive a reasoned answer.
- Service is deployed on Google Cloud and uses at least one managed GCP service.

## Implementation Notes
- Keep hard allergen filtering rule-based; use Gemini for explanation and conversational UX.
- Persist session guardrails (allergies/dislikes) across turns in Redis.
- Add explicit safety disclaimer for severe allergies.
