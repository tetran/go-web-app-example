.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build: ## Build docker image to deploy
	docker build -t kkoichi/gotodo:$(DOCKER_TAG) --target deploy ./

build-local: ## Build docker image for local development
	docker compose build --no-cache

up: ## Do docker compose up
	docker compose up -d

down: ## Do docker compose down
	docker compose down

logs: ## Tail docker compose logs
	docker compose logs -f

ps: ## Check container status
	docker compose ps

test: ## Run tests
	go test -race -shuffle=on ./...

dry-migrate: ## Try migration
	mysqldef -u gotodo -p gotodo -h 127.0.0.1 -P 13306 gotodo --dry-run < ./_tools/mysql/schema.sql

migrate:  ## Execute migration
	mysqldef -u gotodo -p gotodo -h 127.0.0.1 -P 13306 gotodo < ./_tools/mysql/schema.sql

generate: ## Generate codes
	go generate ./...

help: # Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
				awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
