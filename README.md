# create-templ-app

## Stack

- **Go + Gin** — HTTP server, server-side rendering
- **templ** — HTML components compiled to Go
- **Tailwind CSS v4** — compiled via `@tailwindcss/cli`
- **Alpine.js + htmx** — bundled via esbuild
- **esbuild** — JS bundling and minification

## Frontend build

`build/build.js` runs both tools and writes `internal/assets/dist/manifest.json`:

```json
{
  "main.css": "main-a2a25137.css",
  "main.js": "main-6HC7DWYN.js"
}
```

The Go server reads this manifest at startup. `resolver.URL("main.css")` returns the full URL with the hashed filename — `/static/main-a2a25137.css` by default, or wherever `ASSET_BASE_URL` points otherwise.

## Asset serving

By default, Go serves `/static/*` directly from `internal/assets/dist/`. `ASSET_BASE_URL` lets asset URLs point elsewhere instead — e.g. a CDN in front of a bucket you sync `dist/` to — in which case the Go server only ever serves HTML and never touches the static files itself.

Leave `ASSET_BASE_URL` unset to serve assets locally (defaults to `/static`).

## Prerequisites

- Go 1.24+ and Node.js/npm

`templ` and `air` don't need a separate global install — they're pinned as Go tool dependencies (see the `tool (...)` block in `go.mod`, added via `go get -tool`). `go tool templ`/`go tool air` fetch and run the exact pinned version automatically, the same way `require` pins library versions.

## Running locally

```bash
# install frontend deps (once)
npm install

# build frontend assets
npm run build

# generate templ → Go (must run before go build)
go tool templ generate

# start the server
go run ./cmd/server
```

For active development:

```bash
make dev
```

This runs `npm run watch` (rebuilds CSS/JS on change) and `air` (hot-reloads the Go server) in parallel. `air` already runs `templ generate` as part of its own build step (see `.air.toml`), so there's no separate templ watcher process needed.

If you'd rather run things in separate terminals — e.g. to see templ's own watch output — `make assets-watch`, `make templ-watch`, and `make air` are still available individually.

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `ASSETS_DIST_PATH` | `internal/assets/dist` | Path to esbuild output directory |
| `ASSET_BASE_URL` | `/static` | Base URL for asset resolution. Point it wherever assets are actually hosted, if not served locally. |

See `.env.example`.
