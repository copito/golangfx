FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

# Copy Makefile over
COPY Makefile ./Makefile

# Install goose (migration tool), sqlc (repository pattern), proto (protobuf compiler)
RUN make migrate_init sqlc_init proto_init

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# [Optional] Uncomment the next lines to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>
# USER root

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1