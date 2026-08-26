# SPDX-FileCopyrightText: 2025 Deutsche Telekom IT GmbH
#
# SPDX-License-Identifier: Apache-2.0

# Image URL to use all building/pushing image targets
IMG ?= issuer-service:latest

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Format Go code with gofmt.
	gofmt -w .

.PHONY: fmt-check
fmt-check: ## Check whether Go code is formatted.
	@unformatted="$$(gofmt -l .)"; \
	if [[ -n "$${unformatted}" ]]; then \
		printf 'The following Go files are not formatted:\n%s\n' "$${unformatted}"; \
		exit 1; \
	fi

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test-unit
test-unit: ## Run tests.
	go test ./... -race -cover -coverprofile=cover.profile

.PHONY: lint
lint: ## Run golangci-lint linter.
	golangci-lint run

.PHONY: lint-fix
lint-fix: ## Run golangci-lint linter and perform fixes.
	golangci-lint run --fix

.PHONY: lint-config
lint-config: ## Verify golangci-lint linter configuration.
	golangci-lint config verify

.PHONY: reuse
reuse: ## Check REUSE compliance.
	reuse lint

.PHONY: govulncheck
govulncheck: ## Check Go packages for known vulnerabilities.
	govulncheck ./...

.PHONY: check
check: ## Run the source-level CI checks in sequence.
	$(MAKE) fmt-check
	$(MAKE) lint
	$(MAKE) reuse
	$(MAKE) build
	$(MAKE) test-unit
	$(MAKE) govulncheck

.PHONY: hooks
hooks: ## Verify prerequisites and install the opt-in Git hooks.
	@missing=(); \
	for tool in gofmt lefthook reuse gitleaks committed; do \
		if ! command -v "$${tool}" >/dev/null 2>&1; then \
			missing+=("$${tool}"); \
		fi; \
	done; \
	if (( $${#missing[@]} > 0 )); then \
		printf 'Missing hook prerequisites: %s\nSee CONTRIBUTING.md for installation instructions.\n' "$${missing[*]}"; \
		exit 1; \
	fi
	lefthook install

##@ Build

.PHONY: build
build: ## Build issuer-service binary.
	go build -o issuer-service ./cmd/api/main.go

.PHONY: run
run:  ## Run a controller from your host.
	go run ./cmd/main.go

# If you wish to build the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	$(CONTAINER_TOOL) build -t ${IMG} -f Dockerfile.multi-stage .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) push ${IMG}


##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif
