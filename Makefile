# Makefile to help get up and running with Application

PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc
PROTOC_GEN_GRPC_GATEWAY := protoc-gen-grpc-gateway

# point to proto files outside service
PROTO_DIR := idl/proto
OPENAPI_DIR := openapi
PROTO_OUTPUT_DIR := idl_gen

# find all proto files
PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

# Entrypoint service
ENTRYPOINT_PATH := src/cmd/service/main.go
# DB Migrations
DB_MIGRATION_PATH := ./db_migrations

# Build the application
all: build test


########################################
########### PROTO ######################
########################################

proto_init:
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/bufbuild/buf/cmd/buf@v1.47.2 
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	@go install google.golang.org/protobuf
	# @buf config init # generates buf.yaml file

# adding https://github.com/googleapis/googleapis -> annotations.proto / http.proto (via deps)
# More examples for grpc-gateway at: https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/examples/
proto:
	@buf dep update ${PROTO_DIR}
	@buf generate
	@go get google.golang.org/grpc && go get github.com/grpc-ecosystem/grpc-gateway/v2/runtime
	@go mod tidy

clean_proto:
	rm -rf $(PROTO_OUTPUT_DIR)/* && touch $(PROTO_OUTPUT_DIR)/.keep
	rm -rf $(OPENAPI_DIR)/* && touch $(OPENAPI_DIR)/.keep


proto_lint:
	@buf lint ${PROTO_DIR}

########################################
########### Application   ##############
########################################

build:
	@echo "Building..."
	@go build -o main ${ENTRYPOINT_PATH}

# Run the application
run:
	@go run ${ENTRYPOINT_PATH}

# Create DB container
docker_run:
	@if docker compose -f .devcontainer/docker-compose.yaml up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f .devcontainer/docker-compose.yaml up --build; \
	fi

# Shutdown DB container
docker_down:
	@if docker compose -f .devcontainer/docker-compose.yaml down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f .devcontainer/docker-compose.yaml down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./src/... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./src/internal/db -v

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

########################################
####### Migration commands   ###########
########################################

migrate_init:
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

migrate_up:
	@echo "Applying all migrations..."
	@goose -dir ${DB_MIGRATION_PATH} postgres "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" up

migrate_up_one:
	@echo "Applying all migrations..."
	@goose -dir ${DB_MIGRATION_PATH} postgres "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" up-by-one

migrate_down_one:
	@echo "Rolling back the last migration..."
	@goose -dir ${DB_MIGRATION_PATH} postgres "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" down

migrate_create:
	@read -p "Enter migration name: " name; \
	goose create "$$name" sql -dir ./db_migrations

migrate_status:
	@echo "Checking migration status..."
	@goose -dir ${DB_MIGRATION_PATH} postgres "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" status

migrate_reset:
	@echo "Rolling back all migrations..."
	@goose -dir ${DB_MIGRATION_PATH} postgres "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" reset

########################################
####### Migration commands   ###########
########################################


.PHONY: all build run test clean watch docker_run docker_down itest migrate_init migrate_up migrate_up_one migrate_down migrate_create migrate_status migrate_reset install_dev_requirements

