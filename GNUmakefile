VERSION ?= dev
BINARY := terraform-provider-ubuntu

GOLANGCI_LINT := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0
GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck@v1.7.0

SSH_KEY := testdata/.ssh/id_ed25519
COMPOSE := docker compose -f testdata/docker-compose.yml

TEST_HOST ?= 127.0.0.1
TEST_PORT ?= 22022
TEST_USER ?= deploy
UBUNTU_VERSION ?= 24.04

ACC_ENV := TF_ACC=1 UBUNTU_TEST_HOST=$(TEST_HOST) UBUNTU_TEST_PORT=$(TEST_PORT) \
	UBUNTU_TEST_USER=$(TEST_USER) UBUNTU_TEST_PRIVATE_KEY_PATH=$(CURDIR)/$(SSH_KEY)

.PHONY: build test lint govulncheck docs testacc testacc-cover testenv-up testenv-down

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

testacc: $(SSH_KEY)
	$(ACC_ENV) go test ./internal/provider/... -v -count=1 -timeout 30m

testacc-cover: $(SSH_KEY)
	$(ACC_ENV) go test ./... -count=1 -coverpkg=./... -coverprofile=coverage.out -timeout 30m

# ssh ignores a private key that other users can read.
$(SSH_KEY):
	mkdir -p $(dir $(SSH_KEY))
	ssh-keygen -q -t ed25519 -N '' -C ubuntu-provider-acceptance -f $(SSH_KEY)
	chmod 600 $(SSH_KEY)

testenv-up: $(SSH_KEY)
	PUBLIC_KEY="$$(cat $(SSH_KEY).pub)" UBUNTU_TEST_PORT=$(TEST_PORT) UBUNTU_VERSION=$(UBUNTU_VERSION) \
		$(COMPOSE) up -d --build --wait
	@for i in $$(seq 1 60); do \
		ssh -q -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o BatchMode=yes \
			-o ConnectTimeout=2 -i $(SSH_KEY) -p $(TEST_PORT) $(TEST_USER)@$(TEST_HOST) true && exit 0; \
		sleep 1; \
	done; echo "sshd did not come up on $(TEST_HOST):$(TEST_PORT)"; exit 1

testenv-down:
	PUBLIC_KEY=unused $(COMPOSE) down -v --rmi local
