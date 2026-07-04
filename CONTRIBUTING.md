# Contributing

See the README's "Developing this repo" section for setup.

## Two modules

This repo is two independent Go modules: the root module (the app that gets scaffolded), and `cmd/create-templ-app` (the CLI that does the scaffolding, with its own `go.mod` so its code stays out of generated apps. See `cmd/create-templ-app/exclusion_test.go`.

## Scope

This is a template, not a framework. It's meant to stay small and unopinionated enough to fork and modify. Prefer fixes and narrow improvements to the existing stack over adding new dependencies or abstractions.
