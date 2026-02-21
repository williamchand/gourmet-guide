.PHONY: backend-test frontend-test infra-check fmt lint

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
