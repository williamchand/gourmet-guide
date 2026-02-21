#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Deploy the frontend build to a Google Cloud Storage bucket.

Usage:
  scripts/deploy_frontend_gcs.sh --bucket <bucket-name> [--project <gcp-project>] [--skip-build]

Options:
  --bucket      Target GCS bucket name (without gs:// prefix)
  --project     Optional GCP project to use for gcloud/gsutil commands
  --skip-build  Skip npm ci + npm run build (uses existing frontend/dist)
  -h, --help    Show this help output
USAGE
}

BUCKET=""
PROJECT=""
SKIP_BUILD="false"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --bucket)
      BUCKET="${2:-}"
      shift 2
      ;;
    --project)
      PROJECT="${2:-}"
      shift 2
      ;;
    --skip-build)
      SKIP_BUILD="true"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

if [[ -z "$BUCKET" ]]; then
  echo "Error: --bucket is required." >&2
  usage >&2
  exit 1
fi

if ! command -v gsutil >/dev/null 2>&1; then
  echo "Error: gsutil is required but not installed." >&2
  exit 1
fi

if [[ -n "$PROJECT" ]]; then
  gcloud config set project "$PROJECT" >/dev/null
fi

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="$ROOT_DIR/frontend"
DIST_DIR="$FRONTEND_DIR/dist"
BUCKET_URI="gs://$BUCKET"

if [[ "$SKIP_BUILD" != "true" ]]; then
  echo "Building frontend assets..."
  (
    cd "$FRONTEND_DIR"
    if [[ -f package-lock.json ]]; then
      npm ci
    else
      npm install
    fi
    npm run build
  )
fi

if [[ ! -d "$DIST_DIR" ]]; then
  echo "Error: frontend build output not found at $DIST_DIR" >&2
  exit 1
fi

echo "Ensuring bucket exists: $BUCKET_URI"
if ! gsutil ls -b "$BUCKET_URI" >/dev/null 2>&1; then
  if [[ -z "$PROJECT" ]]; then
    echo "Error: bucket does not exist. Provide --project so the script can create it." >&2
    exit 1
  fi
  gsutil mb -p "$PROJECT" -b on "$BUCKET_URI"
fi

echo "Syncing static assets to $BUCKET_URI ..."
gsutil -m rsync -r -d "$DIST_DIR" "$BUCKET_URI"

echo "Applying cache headers..."
gsutil -m setmeta -h "Cache-Control:public,max-age=31536000,immutable" "$BUCKET_URI/assets/**" || true
gsutil setmeta -h "Cache-Control:no-cache,max-age=0,must-revalidate" "$BUCKET_URI/index.html"

if gsutil ls "$BUCKET_URI"/200.html >/dev/null 2>&1; then
  gsutil setmeta -h "Cache-Control:no-cache,max-age=0,must-revalidate" "$BUCKET_URI/200.html"
fi

if gsutil ls "$BUCKET_URI"/404.html >/dev/null 2>&1; then
  gsutil setmeta -h "Cache-Control:no-cache,max-age=0,must-revalidate" "$BUCKET_URI/404.html"
fi

echo "Configuring website index + 404 pages..."
gsutil web set -m index.html -e 404.html "$BUCKET_URI"

echo "Deployment complete."
echo "If bucket-level public access is blocked by policy, front with HTTPS Load Balancer + Cloud CDN."
