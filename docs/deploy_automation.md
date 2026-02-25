# Automated deployment with GitHub Actions

This repository includes `.github/workflows/deploy.yml` for end-to-end automated deployment on version tag pushes (for example `v1.2.0`) or published GitHub releases.

## What gets deployed
1. Backend container is built from `backend/` and pushed to a container registry such as Docker Hub, GitHub Packages, or any registry you configure (tagged with `github.ref_name`, typically your release tag).
2. Terraform applies infrastructure from `infra/` on Google Cloud Platform:
   - Google Cloud Run backend service
   - Firestore and Cloud Storage resources
3. Frontend is built from `frontend/` and synced to a Google Cloud Storage static hosting bucket.

## Required GitHub repository configuration
### Variables
- `GCP_PROJECT_ID`
- `GCP_REGION`
- `DEPLOY_ENVIRONMENT`
- `FRONTEND_GCS_BUCKET`
- `DOCKER_REGISTRY` – the full image name to push (e.g., `docker.io/myorg/gourmet-guide-backend`).

### Secrets
- `DOCKER_USERNAME` (registry user or org)
- `DOCKER_PASSWORD` (registry password or personal access token)
- `GCP_WORKLOAD_IDENTITY_PROVIDER`
- `GCP_SERVICE_ACCOUNT`

## Terraform variables used by CI
- `project_id`
- `environment`
- `region`
- `backend_image`
- `allow_unauthenticated=true`

## Triggering a release deployment
Use semantic version tags for production-style releases:

```bash
git tag v1.0.0
git push origin v1.0.0
```

You can also run by publishing a GitHub Release for an existing tag.
