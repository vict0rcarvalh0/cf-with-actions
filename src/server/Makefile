# Default environment file
ENV_FILE := .env

# Docker container names
CONTAINER_NAME_DEV := postgres_dev
CONTAINER_NAME_PROD := postgres_prod

# Enhanced cleanup function
define docker_cleanup
	echo "Cleaning up Docker containers and volumes..."; \
	docker stop $(1) || true; \
	docker rm $(1) || true; \
	docker compose -f $(2) down -v
endef

# Ensure .env file exists
$(ENV_FILE):
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "Creating .env from .env.example"; \
		cp .env.example $(ENV_FILE); \
	else \
		echo ".env file already exists"; \
	fi

# General target to manage the database container
.PHONY: manage-db-dev
manage-db-dev:
	@if [ -z "$$(docker ps -q -f name=$(CONTAINER_NAME_DEV))" ]; then \
		$(call docker_cleanup,$(CONTAINER_NAME_DEV),compose.dev.yml); \
		docker compose -f compose.dev.yml up --build -d; \
	fi

.PHONY: manage-db-prod
manage-db-prod:
	@if [ -z "$$(docker ps -q -f name=$(CONTAINER_NAME_PROD))" ]; then \
		$(call docker_cleanup,$(CONTAINER_NAME_PROD),compose.prod.yml); \
		docker compose -f compose.prod.yml up --build -d; \
	fi

# Database - runs the Postgres database in a Docker container
.PHONY: database-dev
database-dev: manage-db-dev

.PHONY: database-prod
database-prod: manage-db-prod

# Dev - runs the server in development mode using Air along with the database
.PHONY: dev
dev: $(ENV_FILE) manage-db-dev
	@npx kill-port 8080
	@go mod tidy
	@air

# Prod - runs the application and database by running the compose.prod.yml file
.PHONY: prod
prod: $(ENV_FILE) manage-db-prod
	@npx kill-port 8080
	@docker compose -f compose.prod.yml up --build -d

# Reset - stops and removes the database container
.PHONY: reset-dev
reset-dev:
	@$(call docker_cleanup,$(CONTAINER_NAME_DEV),compose.dev.yml)

.PHONY: reset-prod
reset-prod:
	@$(call docker_cleanup,$(CONTAINER_NAME_PROD),compose.prod.yml)

# Test - runs the tests with database container
.PHONY: test
test: reset-dev $(ENV_FILE) database-dev
	@go test -v ./internal/tests...
	@$(call docker_cleanup,$(CONTAINER_NAME_DEV),compose.dev.yml)

# Deploy - deploys the application
.PHONY: deploy
deploy:
	echo "Deploying application..."; \
	scp -i $(SSH_KEY) ~/.ssh/id_ed25519 ubuntu@$(IP):~/.ssh/id_ed25519; \
	chmod +x deploy.sh; \
	./deploy.sh $(SSH_KEY) $(IP) $(CLONE_SSH_KEY) || { echo "Deployment failed"; exit 1; }

# Build - builds the application
.PHONY: build
build:
	go build -o bin/app cmd/app/main.go

# Clean - removes temporary files/folders, clean up the database container and removes the binary
.PHONY: clean
clean: reset-dev reset-prod
	@rm -rf bin
	@rm -rf .air
	@rm -rf tmp
	@rm -rf go.sum
	@rm -rf vendor
	@rm -rf .DS_Store
	@rm -rf .vscode
	@rm -rf .idea
	@docker compose -f compose.dev.yml down -v
	@docker compose -f compose.prod.yml down -v