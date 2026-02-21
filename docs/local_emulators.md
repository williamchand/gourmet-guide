# Local Firestore + Cloud Storage Emulators

This project can run local emulators for Firestore and Cloud Storage so seed generation and publish flows can be exercised without deploying to GCP.

## Start emulators

```bash
docker compose -f docker-compose.dev.yml up -d
```

Services started:
- Firestore emulator on `localhost:8081`.
- GCS emulator (fake-gcs-server) on `localhost:4443`.
- A one-shot bucket bootstrap that creates `gourmet-guide-seed-local`.

## Environment variables for local runs

```bash
export GOOGLE_CLOUD_PROJECT=gourmet-guide-local
export FIRESTORE_EMULATOR_HOST=localhost:8081
export STORAGE_EMULATOR_HOST=http://localhost:4443
export GCS_BUCKET=gourmet-guide-seed-local
export GOOGLE_API_KEY=replace-with-your-key
```

> `GOOGLE_API_KEY` is still needed for Gemini generation, but Firestore and Storage writes are sent to local emulators when `FIRESTORE_EMULATOR_HOST` and `STORAGE_EMULATOR_HOST` are set.

## Run seed generator against emulators

```bash
cd backend
go run -tags gcp ./cmd/seeddata --count 10 --out seed/output
```

## Quick smoke checks

List firestore docs:

```bash
curl -s "http://localhost:8081/v1/projects/gourmet-guide-local/databases/(default)/documents/restaurants" | jq .
```

List storage buckets:

```bash
curl -s "http://localhost:4443/storage/v1/b" | jq .
```

## Stop emulators

```bash
docker compose -f docker-compose.dev.yml down
```
