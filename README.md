# Project GoBase Runner K8s

This project is supposed to be a proof of concept of building a clean golang service with the best architecture possible to ensure
that it can be reused in the future by other projects. For simplicity sake it will be a microservice that interacts with databases (postgres and timescale), kafka (redpandas), GRPC, Kubernetes and more.

## Prerequisites

1. Golang Version Manager: recommend installing `asdf` or `gvm` (go versiion manager) to manage multiple golang version, but if you prefer to install go directly with the version specified in the [go.mod file](./go.mod)

ASDF:

```bash
asdf plugin-add golang https://github.com/kennyp/asdf-golang.git
asdf install golang 1.23.5
asdf global golang 1.23.5
```

2. Just: recommend installing just in this project, which is an alternative to Makefile - which is a handy way to save and run project-specific commands. It can be downloaded directly from [Source](https://github.com/casey/just) or simply install using package manager:

```bash
sudo apt install just
```

3. Other internal dependencies: these can be downloaded via the just commands that end with `_init` - such as `proto_init`. All of the commands can be found under the [justfile](./justfile) and should be installed as needed based on what changes you are making (database migration, proto changes, sql generation for repositories, etc)

- Proto: using [buf](https://buf.build/) as the tool of choice
- SQL Generation: using [sqlc](https://sqlc.dev/) as the tool of choice
- Database Migrations: using [goose](https://pressly.github.io/goose)
- Hot Reloading: using [air](https://github.com/air-verse/air)

## Structure

```
root/
├── .devcontainer/             # Manages all local development scripts (devcontainers)
├── .vscode/                   # Vscode specific configs to help developer experience (i.e. local debugging)
├── db_migrations/             # Keeps all migrations for goose to migrate
├── deploy/                    # Keeps all deployment files for Kubernetes
├── idl/                       # Keeps all Interface definition Languages like Proto, Thift, etc
├── idl_gen/                   # Keeps all the generated code from the idl/ folder
├── openapi/                   # Keeps all the generated openapi specs
├── config
    └── base.yaml        # Base configuration file loaded by Viper (make sure to add to entity/config, if changed)
├── src/
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
    │   ├── proto/                 # Compiled protobuf files
    ├── pkg/                       # Reusable packages (utility functions, helpers)
    │   ├── logger/                # fx-compatible logging module
    │   ├── middleware/            # fx-compatible middlewares for gRPC and HTTP
    │   └── utils/                 # General utility functions
├── buf.gen.yaml            # Actually handles the creation of auto-generated code (grpc, openapi, grpc-gateway)
├── buf.yaml                # Handles the dependency management of idl/proto external files
├── buf.lock                # Handles the lock for these external dependencies
├── go.sum
└── go.mod                  # Go module file

```

## Local Development

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.
There are two ways to test this system:

1. The first is via devcontainers (recommended) as it installs any necessary packages for you in a container and spins up any dependencies in the same manner. Devcontainers are also an excellent way to help keep development consistent as it allows to suggest extensions in VSCode as soon as you start it.

1. The second way, is to simply run locally using the make commands provided (less recommended). This way, you will have to have some packages installed globally that take care of proto generation, auto-reload (for faster development), database migrations and more, but all of these are provided via the Makefile.

## MakeFile

### Application Specific

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

### Protobuf / GRPC specific

This command prepares protobuf for your system, but installing all global dependencies.
This will not be required if using devcontainer

```bash
make proto_init
```

This command generates all idl_gen files from the proto files in idl/,
by running `buf`. [Read more](https://buf.build/blog/buf-cli-next-generation)

```bash
make proto
```

This command removes all [idl_gen/](./idl_gen/) (auto generated code) and [openapi/](./openapi/) (auto generated specs).
This is useful if you want to start from a clean slate before running `make proto`

```bash
make clean_proto
```

### Docker - Local development

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

### Database Migrations

This command prepares database migration tooling for your system, by installing all global dependencies.
This will not be required if using devcontainer. [goose](https://github.com/pressly/goose)

```bash
make migrate_init
```

The command `migrate_up` ensures that you run a migration up to the latest
offset based on what is found on the database. This will be running during CI (continuous integration)
to make sure it is a breaking change.

```bash
make migrate_up
```

The command `migrate_up_one` ensures that you run a migration only for the next
offset based on what is found on the database. This helps for testing applying
migrations one at a time.

```bash
make migrate_up_one
```

The command `migrate_down_one` ensures that you run a migration down by one - rolling back
what has just been done. This helps when there are changes that have failed or reset cases (bad deployments)

```bash
make migrate_down_one
```

The command `migrate_create` ensures that you create a new migration in a consistent manner
in the correct folder and with the correct format so that it picks up from the last one.

```bash
make migrate_create
```

The command `migrate_status` gives you the migration status with information when were the
previous migrations applied and until which point have you applied those.

_Example_:

```c
$   Applied At                  Migration
$   =======================================
$   Sun Jan  6 11:25:03 2013 -- 001_basics.sql
$   Sun Jan  6 11:25:03 2013 -- 002_next.sql
$   Pending                  -- 003_and_again.go
```

```bash
make migrate_status
```

The command `migrate_reset` has to be used with a lot of caution.
It will reset all migrations by rolling back all migrations.
This should only be used in **Local Development**!

```bash
make migrate_reset
```
