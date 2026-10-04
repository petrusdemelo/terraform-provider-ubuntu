VERSION ?= dev
BINARY := terraform-provider-ubuntu

GOLANGCI_LINT := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0
GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck@v1.7.0

.PHONY: build test lint govulncheck docs

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "-X main.version=$(VERSION)" -o $(BINARY) .

test:
	go test ./...

lint:
	go run $(GOLANGCI_LINT) run

govulncheck:
	go run $(GOVULNCHECK) ./...

docs:
	go tool tfplugindocs generate --provider-name ubuntu
