######################################################################################################
#                                                                                                    #
#                                   EUI-64 Calculator Makefile                                       #
#                                                                                                    #
######################################################################################################

# Binary and module configuration.
BINARY_NAME := eui64-calculator
MODULE       := github.com/nicholas-fedor/eui64-calculator
ENTRYPOINT   := ./cmd/server/main.go

# Build output directory.
DIST_DIR := dist

# Build metadata injected via ldflags.
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# ldflags for version injection.
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Go and tool configuration.
CGO_ENABLED := 0
GOFLAGS     := -trimpath

# golangci-lint configuration path.
GOLANGCI_LINT_CONFIG := build/golangci-lint/golangci-lint.yaml

# GoReleaser configuration path.
GORELEASER_CONFIG := build/goreleaser/goreleaser.yaml

# Docker configuration.
DOCKER_IMAGE  := $(BINARY_NAME)
DOCKER_TAG    ?= latest

# Test configuration.
TEST_FLAGS := -v -count=1
COVER_DIR  := coverage

######################################################################################################
# Phony targets
######################################################################################################
.PHONY: all check clean cover bench docker-build docker-run fmt generate help \
        lint mod-tidy release run test test-race test-coverage vet

# Default target.
all: check ## Run all validation checks (lint, vet, test)

######################################################################################################
# Help
######################################################################################################

# Self-documenting Makefile: extracts targets and their ## comments.
help: ## Display this help message
	@printf '\nUsage:\n  make \033[36m<target>\033[0m\n\n'
	@printf '\033[1mAvailable targets:\033[0m\n'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@printf '\n'

######################################################################################################
# Code Generation
######################################################################################################

generate: ## Run go generate to produce templ and other generated code
	go generate ./...

######################################################################################################
# Linting and Validation
######################################################################################################

lint: ## Run golangci-lint with project configuration
	golangci-lint run --config $(GOLANGCI_LINT_CONFIG) ./...

vet: ## Run go vet for additional static analysis
	go vet ./...

check: lint vet test ## Run lint, vet, and test (full CI check)

######################################################################################################
# Formatting
######################################################################################################

fmt: ## Format code and organize imports with golangci-lint
	golangci-lint fmt --config $(GOLANGCI_LINT_CONFIG) ./...

######################################################################################################
# Testing
######################################################################################################

test: generate ## Run all tests
	go test $(TEST_FLAGS) ./...

test-race: generate ## Run all tests with the race detector enabled
	go test $(TEST_FLAGS) -race ./...

test-coverage: generate ## Run all tests with coverage reporting
	@mkdir -p $(COVER_DIR)
	go test $(TEST_FLAGS) -race -coverprofile=$(COVER_DIR)/coverage.out -covermode=atomic ./...
	go tool cover -html=$(COVER_DIR)/coverage.out -o $(COVER_DIR)/coverage.html
	@printf '\nCoverage report: $(COVER_DIR)/coverage.html\n'

cover: test-coverage ## Alias for test-coverage

bench: generate ## Run all benchmark tests
	go test -bench=. -benchmem ./...

######################################################################################################
# Build
######################################################################################################

run: generate ## Run the server locally
	CGO_ENABLED=$(CGO_ENABLED) go run $(GOFLAGS) -ldflags "$(LDFLAGS)" $(ENTRYPOINT)

######################################################################################################
# Module Management
######################################################################################################

mod-tidy: ## Tidy and verify go.mod dependencies
	go mod tidy
	go mod verify

######################################################################################################
# Docker
######################################################################################################

docker-build: ## Build the Docker image
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 \
		go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME) $(ENTRYPOINT)
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f build/docker/Dockerfile $(DIST_DIR)

docker-run: docker-build ## Build and run the Docker container
	docker run --rm -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

######################################################################################################
# Release
######################################################################################################

release: ## Create a release build with GoReleaser
	goreleaser release --config $(GORELEASER_CONFIG)

######################################################################################################
# Cleanup
######################################################################################################

clean: ## Remove build artifacts and generated files
	rm -rf $(DIST_DIR) $(COVER_DIR)
	go clean -testcache

######################################################################################################
