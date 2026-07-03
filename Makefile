.PHONY: help
help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-14s\033[0m %s\n", $$1, $$2}'

.PHONY: dev
dev: ## Run the frontend watcher and air together for local development
	@$(MAKE) -j2 assets-watch air

.PHONY: run
run: ## Run the server without hot-reloading
	go run ./cmd/server

.PHONY: build
build: assets templ ## Full production build (assets + templ + go build)
	go build -o bin/server ./cmd/server

.PHONY: test
test: ## Run unit tests
	go test ./...

.PHONY: assets
assets: ## One-shot frontend build
	npm run build

.PHONY: assets-watch
assets-watch: ## Frontend build in watch mode (for dev)
	npm run watch

.PHONY: templ
templ: ## Generate templ -> Go
	go tool templ generate

.PHONY: templ-watch
templ-watch: ## Generate templ -> Go in watch mode (for dev)
	go tool templ generate --watch

.PHONY: air
air: ## Go hot-reloading via air (also runs templ generate on every rebuild)
	go tool air

.PHONY: tidy
tidy: ## Tidy Go modules and install npm dependencies
	go mod tidy
	npm install
