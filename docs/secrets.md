# Secrets Management

## Local Development
- Copy `.env.example` to `.env` in repository root.
- Use short-lived API keys for local development only.
- Never commit `.env` files or plain-text credentials.

## Cloud Environments (Google Cloud)
- Prefer Workload Identity + service account auth for Vertex AI, Firestore, and Cloud Storage.
- Store only non-identity secrets in Secret Manager (for example, third-party API credentials).
- Bind Cloud Run service account with least privilege for only required services.
- Rotate secrets on a fixed schedule and after incidents.

## Required Secrets
- `GOOGLE_API_KEY` is optional for local non-Vertex testing.
