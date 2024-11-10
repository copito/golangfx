# Project GoBase Runner K8s

This project is supposed to be a proof of concept of building a clean golang service with the best architecture possible to ensure
that it can be reused in the future by other projects. For simplicity sake it will be a microservice that interacts with databases (postgres and timescale), kafka (redpandas), GRPC, Kubernetes and more.

## Structure

```
root/
├── .devcontainer/             # Manages all local development scripts (devcontainers)
├── db_migrations/             # Keeps all migrations for goose to migrate
├── deploy/                    # Keeps all deployment files for Kubernetes
├── cmd/
│   └── main/                  # Single entry point for all components
│       └── main.go
├── internal/
│   ├── app/                   # Core application logic
│   │   ├── controllers/       # Controller layer for handling HTTP/gRPC requests
│   │   ├── entities/          # Domain entities or data structures
│   │   ├── gateways/          # Interfaces for external integrations (DB, API calls, etc.)
│   │   ├── models/            # Database models or data representations
│   │   └── modules/           # fx modules that group dependencies by component
│   │       ├── web/           # fx module for web and gRPC (gateway) dependencies
│   │       ├── grpc/          # fx module for gRPC dependencies
│   │       ├── kafka/         # fx module for Kafka consumers and producers
│   │       ├── services/      # fx module for business services
│   │       └── config/        # fx module for configuration setup
│   ├── db/                    # Database modules, used as fx components
│       ├── postgres/          # Postgres connection and setup
│       └── timeseries/        # Timeseries database connection and setup
│   ├── config/                # Configuration with Viper
│       ├── config.go
│       └── config.yaml        # Configuration file loaded by Viper
│   ├── proto/                 # Compiled protobuf files
├── api/                       # API schemas (e.g., protobuf, OpenAPI)
│   ├── grpc/                  # gRPC protobuf files
│   └── http/                  # OpenAPI or other HTTP API specs
├── pkg/                       # Reusable packages (utility functions, helpers)
│   ├── logger/                # fx-compatible logging module
│   ├── middleware/            # fx-compatible middlewares for gRPC and HTTP
│   └── utils/                 # General utility functions
└── go.mod                     # Go module file

```

## Local Development

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```
