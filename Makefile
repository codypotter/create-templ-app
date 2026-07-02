.PHONY: dev build assets templ tidy air

# Run the frontend watcher and air together. air regenerates templ output
# and rebuilds/restarts the Go server on its own, so that's the only other
# process needed. Uses make's own -j job control for startup/shutdown
# instead of hand-rolled process management.
dev:
	@$(MAKE) -j2 assets-watch air

run:
	go run ./cmd/server

# Go hot-reloading via air (also runs templ generate on every rebuild)
air:
	air

# One-shot frontend build
assets:
	npm run build

# Frontend in watch mode (for dev)
assets-watch:
	npm run watch

# Generate templ → Go
templ:
	templ generate

# Templ in watch mode (for dev)
templ-watch:
	templ generate --watch

# Full production build
build: assets templ
	go build -o bin/server ./cmd/server

tidy:
	go mod tidy
	npm install
