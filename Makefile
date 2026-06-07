.PHONY: dev build test migrate clean frontend help

# =============================================================================
# Game Platform — Docker + Local Development
# =============================================================================

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

dev: ## Run full stack via docker-compose (build + up)
	docker compose up --build

dev-detach: ## Run full stack detached
	docker compose up --build -d

build: ## Build docker images
	docker compose build

down: ## Stop containers (keep volumes)
	docker compose down

clean: ## Stop containers and remove volumes
	docker compose down -v

test: ## Run Go tests
	go test ./...

test-cover: ## Run Go tests with coverage report
	go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

migrate: ## Run migrations against running database
	./migrations/run.sh

logs: ## Tail container logs
	docker compose logs -f game-platform

psql: ## Open psql shell to the database container
	docker compose exec db psql -U game -d game_platform

frontend: ## Run frontend dev server (Vite)
	cd web && npm run dev

frontend-install: ## Install frontend dependencies
	cd web && npm install
