# create-templ-app

Go + Gin + templ + Tailwind CSS v4 + htmx + Alpine.js, scaffolded with a single command.

## Usage

```bash
go run github.com/codypotter/create-templ-app/cmd/create-templ-app@latest github.com/you/your-app
cd your-app
make tidy
make dev
```

This copies the repo itself into `./your-app` (or a second argument, if given) and rewrites the module path throughout. You end up with a normal, standalone Go module and no dependency on this repo. Open http://localhost:8080; the two example routes under `internal/views` exist to show the htmx/Alpine wiring and are meant to be deleted once you don't need them.

## Why this stack

Server-rendered HTML with htmx and Alpine.js handling interactivity, instead of a client-side framework. `templ` compiles views to typed Go rather than parsing templates at runtime, and Tailwind/esbuild are wired up with content-hashed, cacheable output. The goal is to skip the SPA tax for the large majority of apps that don't need it.

## Stack

- **Go + Gin** — HTTP server, server-side rendering
- **templ** — HTML components compiled to Go
- **Tailwind CSS v4** — compiled via `@tailwindcss/cli`
- **Alpine.js + htmx** — bundled via esbuild

## Prerequisites

- Go 1.25+
- Node.js/npm

## Developing this repo

```bash
npm install
make dev
```

Open http://localhost:8080. Rebuilds CSS/JS and hot-reloads the Go server on every change. Run `make help` to see all available targets.

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
