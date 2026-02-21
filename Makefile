.PHONY: backend-test frontend-test infra-check fmt lint deploy-frontend-gcs

backend-test:
	cd backend && go test ./...

frontend-test:
	cd frontend && npm test

infra-check:
	cd infra && terraform fmt -check && terraform init -backend=false && terraform validate

fmt:
	cd backend && gofmt -w ./cmd ./internal
	cd infra && terraform fmt

lint:
	cd frontend && npm run lint && npm run format


deploy-frontend-gcs:
	@if [ -z "$$BUCKET" ]; then echo "Usage: make deploy-frontend-gcs BUCKET=<bucket> [PROJECT=<project>]"; exit 1; fi
	./scripts/deploy_frontend_gcs.sh --bucket "$$BUCKET" $(if $(PROJECT),--project "$(PROJECT)",)
