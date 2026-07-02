.PHONY: dev build assets templ tidy

# Start frontend watcher, templ watcher, and Go server together
dev:
	@echo "run 'make assets-watch', 'make templ-watch', and 'make run' in separate terminals"
	@echo "or use a tool like overmind/foreman with a Procfile"

run:
	go run ./cmd/server

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
