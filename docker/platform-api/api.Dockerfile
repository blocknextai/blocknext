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

ARG BINARY_NAME
ARG TARGETARCH
ENV GOARCH=${TARGETARCH}

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd apps/platform-api && \
    go build -ldflags="-s -w" -trimpath -o /build/app ./cmd/${BINARY_NAME}/main.go

FROM gcr.io/distroless/static-debian13:nonroot AS production

ARG BINARY_NAME
ARG PORT=3000

WORKDIR /app

COPY --from=build /build/app .
COPY --from=build /build/apps/platform-api/prompts ./prompts

EXPOSE ${PORT}

USER nonroot

ENTRYPOINT ["/app/app"]

LABEL org.opencontainers.image.title="${BINARY_NAME}" \
      org.opencontainers.image.description="BlockNext platform HTTP server (${BINARY_NAME})" \
      org.opencontainers.image.vendor="BlockNext AI" \
      org.opencontainers.image.source="https://github.com/blocknextai/blocknext" \
      org.opencontainers.image.url="https://www.blocknext.ai"
