# create-templ-app

Go + Gin + templ + Tailwind CSS v4 + htmx + Alpine.js, scaffolded with a single command.

## Stack

- **Go + Gin** — HTTP server, server-side rendering
- **templ** — HTML components compiled to Go
- **Tailwind CSS v4** — compiled via `@tailwindcss/cli`
- **Alpine.js + htmx** — bundled via esbuild

## Prerequisites

- Go 1.24+
- Node.js/npm

## Getting started

```bash
npm install
make dev
```

Open http://localhost:8080. Rebuilds CSS/JS and hot-reloads the Go server on every change.

## Build

```bash
make build
```

Outputs `bin/server`.

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `ASSETS_DIST_PATH` | `internal/assets/dist` | Path to esbuild output directory |
| `ASSET_BASE_URL` | `/static` | Base URL for resolved asset URLs |

See `.env.example`.
