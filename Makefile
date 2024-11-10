# Makefile to help get up and running with Application

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main cmd/service/main.go

# Run the application
run:
	@go run cmd/service/main.go

# Create DB container
docker-run:
	@if docker compose -f .devcontainer/docker-compose.yaml up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f .devcontainer/docker-compose.yaml up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose -f .devcontainer/docker-compose.yaml down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f .devcontainer/docker-compose.yaml down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Migration commands
migrate-init:
	@if command -v goose > /dev/null; then \
            goose --version; \
            echo "Goose initialized...";\
        else \
            read -p "Go's 'goose' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/pressly/goose/v3/cmd/goose@latest; \
                goose --version; \
                echo "Goose installed and initialized...";\
            else \
                echo "You chose not to install goose. Exiting..."; \
                exit 1; \
            fi; \
        fi


migrate-up:
	@echo "Applying all migrations..."
	@goose -dir ./db_migrations postgres "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" up

migrate-up-one:
	@echo "Applying all migrations..."
	@goose -dir ./db_migrations postgres "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" up-by-one

migrate-down-one:
	@echo "Rolling back the last migration..."
	@goose -dir ./db_migrations postgres "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" down

migrate-create:
	@read -p "Enter migration name: " name; \
	goose create "$$name" sql -dir ./db_migrations

migrate-status:
	@echo "Checking migration status..."
	@goose -dir ./db_migrations postgres "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" status

migrate-reset:
	@echo "Rolling back all migrations..."
	@goose -dir ./db_migrations postgres "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" reset

.PHONY: all build run test clean watch docker-run docker-down itest migrate-init migrate-up migrate-up-one migrate-down migrate-create migrate-status migrate-reset

