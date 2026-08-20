FROM oven/bun:1.4.0-alpine AS build
WORKDIR /repo

COPY package.json bun.lock ./
COPY apps/platform/package.json ./apps/platform/
RUN --mount=type=cache,target=/root/.bun/install/cache \
    bun install --frozen-lockfile

COPY apps/platform ./apps/platform
RUN cd apps/platform && bun run build

FROM nginxinc/nginx-unprivileged:1.31.3-alpine AS production

ENV PORT=4000

COPY docker/platform/nginx.conf.template /etc/nginx/templates/default.conf.template
COPY --chmod=755 docker/platform/env.sh /docker-entrypoint.d/40-env.sh
COPY --from=build --chown=nginx:nginx /repo/apps/platform/dist /usr/share/nginx/html

EXPOSE 4000

CMD ["nginx", "-g", "daemon off;"]

LABEL org.opencontainers.image.title="platform" \
      org.opencontainers.image.description="BlockNext Platform UI" \
      org.opencontainers.image.vendor="BlockNext AI" \
      org.opencontainers.image.source="https://github.com/blocknextai/blocknext" \
      org.opencontainers.image.url="https://www.blocknext.ai"
