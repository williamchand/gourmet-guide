# Frontend deployment to Google Cloud Storage (GCS)

This project can be deployed as a static site using a GCS bucket.

## Prerequisites
- Google Cloud project with billing enabled
- `gcloud` + `gsutil` installed and authenticated
- Permissions:
  - `storage.buckets.create` (only if creating bucket)
  - `storage.objects.create`
  - `storage.objects.delete`
  - `storage.objects.update`

## One-time setup
```bash
gcloud auth login
gcloud config set project <your-project-id>
```

## Deploy command
From repository root:

```bash
scripts/deploy_frontend_gcs.sh --bucket <your-bucket-name> --project <your-project-id>
```

What it does:
1. Builds frontend assets (`npm ci && npm run build` when lockfile exists, otherwise `npm install && npm run build`)
2. Creates the bucket if it does not exist
3. Syncs `frontend/dist` to the bucket
4. Sets cache headers:
   - long cache for `assets/**`
   - no-cache for `index.html`
5. Configures static website index/404 behavior

## SPA routing note
GCS static website hosting returns only the configured 404 page for unknown routes.

For React Router deep links in production, use one of these options:
- Add a `404.html` that bootstraps the app the same way as `index.html`
- Or put the bucket behind External HTTPS Load Balancer + URL rewrite + Cloud CDN

## Optional: skip build
If you already built assets:

```bash
scripts/deploy_frontend_gcs.sh --bucket <your-bucket-name> --skip-build
```
