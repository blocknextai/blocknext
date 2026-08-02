# platform

The BlockNext web application — a React SPA where users build, run, and monitor AI-powered workflows on a visual canvas.

## 🚀 Quick Start

### Prerequisites

- [Bun](https://bun.sh) (see the root `package.json` for the pinned version)
- A running backend — start the stack from the monorepo root (see the root `README.md`)

### Development

Install dependencies and generate the root `.env` once from the monorepo root:

```bash
make setup
```

Then start the dev server from this directory:

```bash
bun run dev
```

Vite reads environment variables from the monorepo root `.env` (`envDir` points there); only `VITE_*`-prefixed keys are exposed to the client.

### Production

The production image (`docker/platform/nginx.Dockerfile`) serves the static build with nginx. Environment values are injected at **runtime**, not build time: an entrypoint script writes `env.js` (`window.__ENV__`) from the container environment, and `src/lib/config.ts` falls back to `import.meta.env` during development.

## 🏗️ Structure

```
src/
├── components/   # Shared UI components
├── features/     # Feature modules (workflow editor, credentials, MCP, ...)
├── hooks/        # Shared hooks
├── layouts/      # Page layouts
├── lib/          # Config, API client, utilities
├── pages/        # Route-level pages
├── stores/       # Zustand stores
└── routes.tsx    # Route definitions
```

Each feature pairs a service (API calls) with a `use-<feature>` hook that owns its state — pages stay stateless.

## 🛠️ Key Technologies

- **Framework**: React 19 + TypeScript, built with Vite
- **Routing**: React Router 8
- **State**: Zustand
- **Styling**: Tailwind CSS 4 + Radix UI primitives
- **Workflow Canvas**: `@xyflow/react` (React Flow)
- **i18n**: i18next
