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

The Go server reads this manifest at startup. `resolver.URL("main.css")` returns the full URL with the hashed filename. In dev that's `/static/main-a2a25137.css`; in prod it's `https://assets.com/main-a2a25137.css`.

## Asset serving

Locally, Go serves `/static/*` directly from `internal/assets/dist/`. In production, static assets are synced to S3 and served via CloudFront. The Go server only serves HTML — it never touches static files in prod.

`ASSET_BASE_URL` controls the base URL for asset resolution. Leave it unset in dev (defaults to `/static`); set it to the CDN URL in prod.

## Running locally

```bash
# install frontend deps (once)
npm install

# build frontend assets
npm run build

# generate templ → Go (must run before go build)
templ generate

# start the server
go run ./cmd/server
```

For active development, run these in parallel:

```bash
npm run watch        # rebuilds CSS/JS on change
templ generate --watch  # regenerates Go code on .templ changes
air                  # hot-reloads the Go server
```

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `ASSETS_DIST_PATH` | `internal/assets/dist` | Path to esbuild output directory |
| `ASSET_BASE_URL` | `/static` | Base URL for asset resolution. Set to CDN URL in prod. |

See `.env.example`.
