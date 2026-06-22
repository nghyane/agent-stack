.PHONY: help dev db-up db-down sqlc migrate-new lint test build openapi tidy

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

dev: ## Run API with hot reload (air)
	air

db-up: ## Start local Postgres
	docker compose up -d

db-down: ## Stop local Postgres
	docker compose down

sqlc: ## Regenerate Go code from SQL
	sqlc generate

migrate-new: ## Create a migration: make migrate-new name=add_posts
	goose -dir db/migrations create $(name) sql

lint: ## Run golangci-lint
	golangci-lint run ./...

test: ## Run tests
	go test ./... -race

build: ## Build the binary
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/api ./cmd/api

openapi: ## Export OpenAPI spec to openapi.yaml (for frontend codegen)
	go run ./cmd/api openapi > openapi.yaml

tidy: ## Tidy go modules
	go mod tidy
