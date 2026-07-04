## What & why

<!-- What does this change, and why? -->

## Testing

<!-- How did you verify this? e.g. `make test`, `make build`, manually scaffolding a test app -->

## Checklist

- [ ] `go test ./...` passes (both root module and `cmd/create-templ-app`, if touched)
- [ ] `go tool templ generate` has no diff, if `.templ` files changed
- [ ] `npm run build` succeeds, if frontend assets changed
