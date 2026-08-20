FROM golang:1.27.0-alpine AS base

WORKDIR /build

COPY apps/file-gateway-api/go.mod apps/file-gateway-api/go.sum ./apps/file-gateway-api/
COPY packages/go-packages/go.mod packages/go-packages/go.sum ./packages/go-packages/

RUN --mount=type=cache,target=/go/pkg/mod \
    cd apps/file-gateway-api && go mod download

COPY packages/go-packages ./packages/go-packages

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

FROM base AS build

COPY apps/file-gateway-api ./apps/file-gateway-api

ARG TARGETARCH
ENV GOARCH=${TARGETARCH}

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd apps/file-gateway-api && \
    go build -ldflags="-s -w" -trimpath -o /build/file-gateway-api ./cmd/file-gateway-api/main.go

RUN mkdir -p /uploads/public /uploads/private && chown -R 65532:65532 /uploads

FROM gcr.io/distroless/static-debian13:nonroot AS production

WORKDIR /app

COPY --from=build /build/file-gateway-api .
COPY --from=build --chown=nonroot:nonroot /uploads /data/uploads

EXPOSE 3300

USER nonroot

ENTRYPOINT ["/app/file-gateway-api"]

LABEL org.opencontainers.image.title="file-gateway-api" \
      org.opencontainers.image.description="BlockNext File Gateway API HTTP server" \
      org.opencontainers.image.vendor="BlockNext AI" \
      org.opencontainers.image.source="https://github.com/blocknextai/blocknext" \
      org.opencontainers.image.url="https://www.blocknext.ai"
