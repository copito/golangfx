# Stage1: Build the go Application
FROM docker.io/golang:1.26.4-trixie AS builder

ENV GOPRIVATE=gitea.example.com,gitea* \
    GOINSECURE=gitea/* \
    GIT_SSL_NO_VERIFY=true \
    GOPROXY=https://artifacts.example/met/repository/go-int,https://proxy.golang.org,direct \
    GO111MODULE=on \
    GIT_SSH_COMMAND="ssj -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"

# Install librdkafka, git, openssh-client and necesary build tools
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    librdkafka-dev \
    git \
    openssh-client \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Copy the rest of the application source code
COPY . .

# Sets up temporary GIT_TOKEN (if needed) to run the go get
RUN --mount=type-secret,id=GIT_TOKEN \
    GIT_TOKEN=$(cat /run/secrets/GIT_TOKEN) && \
    git config --global url."https://x-access-token:${GIT_TOKEN}@gitea.example.com".insteadOf "https://gitea.example.com/" && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build --buildvcs=false -o bin/service ./src/cmd/service

# Stage 2: Creat the final image
FROM docker.io/debian:trixie-slim

# Install the necessary runtime deps
RUN apt-get update && apt-get install -y librdkafka1 ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the binary from the builder stage
COPY --from=builder /app/bin/service /service

# Copy the OpenAPI/Swagger files
COPY --from=builder /app/openapi /openapi/

# Expose the port for your app
EXPOSE 50051 5005

# Set entrypoint
ENTRYPOINT [ "/service" ]
