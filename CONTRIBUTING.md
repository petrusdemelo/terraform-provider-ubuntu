# Contributing

## Requirements

- Go, at the version in `go.mod` (`GOTOOLCHAIN=auto` fetches it)
- Terraform 1.8 or later
- Docker, for the acceptance tests

## Checks

Run these before opening a pull request:

```sh
go build ./... && go vet ./...
make test
make lint
make govulncheck
```

## Commits

Commits follow [Conventional Commits](https://www.conventionalcommits.org/):
`feat:`, `fix:`, `docs:`, `test:`, `ci:`, `chore:`. Every user-visible change
adds a line under `## Unreleased` in `CHANGELOG.md`.
