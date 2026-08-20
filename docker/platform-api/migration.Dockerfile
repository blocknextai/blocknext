FROM golang:1.27.0-alpine AS base

WORKDIR /build

COPY apps/platform-api/go.mod apps/platform-api/go.sum ./apps/platform-api/
COPY packages/go-packages/go.mod packages/go-packages/go.sum ./packages/go-packages/

RUN --mount=type=cache,target=/go/pkg/mod \
    cd apps/platform-api && go mod download

COPY packages/go-packages ./packages/go-packages

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

FROM base AS build

COPY apps/platform-api ./apps/platform-api

ARG TARGETARCH
ENV GOARCH=${TARGETARCH}

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd apps/platform-api && \
    go build -ldflags="-s -w" -trimpath -o /build/platform-api-migration ./cmd/platform-api-migration/main.go

FROM gcr.io/distroless/static-debian13:nonroot AS production

WORKDIR /app

COPY --from=build /build/platform-api-migration .

COPY --from=build /build/apps/platform-api/internal/common/infrastructure/database/migrations /app/migrations/common/
COPY --from=build /build/apps/platform-api/internal/account/infrastructure/database/migrations /app/migrations/account/
COPY --from=build /build/apps/platform-api/internal/organizations/infrastructure/database/migrations /app/migrations/organizations/
COPY --from=build /build/apps/platform-api/internal/executions/infrastructure/database/migrations /app/migrations/executions/
COPY --from=build /build/apps/platform-api/internal/triggers/infrastructure/database/migrations /app/migrations/triggers/
COPY --from=build /build/apps/platform-api/internal/workflows/infrastructure/database/migrations /app/migrations/workflows/
COPY --from=build /build/apps/platform-api/internal/credentials/infrastructure/database/migrations /app/migrations/credentials/
COPY --from=build /build/apps/platform-api/internal/apikeys/infrastructure/database/migrations /app/migrations/apikeys/
COPY --from=build /build/apps/platform-api/internal/notifications/infrastructure/database/migrations /app/migrations/notifications/
COPY --from=build /build/apps/platform-api/internal/eventbus/infrastructure/database/migrations /app/migrations/eventbus/

USER nonroot

ENTRYPOINT ["/app/platform-api-migration"]

LABEL org.opencontainers.image.title="platform-api-migration" \
      org.opencontainers.image.description="BlockNext platform database migration runner" \
      org.opencontainers.image.vendor="BlockNext AI" \
      org.opencontainers.image.source="https://github.com/blocknextai/blocknext" \
      org.opencontainers.image.url="https://www.blocknext.ai"
