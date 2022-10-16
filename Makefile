.DEFAULT_GOAL := help

VERSION := $(shell git describe --tags --always --dirty --match "v[0-9]+(\.[0-9]+)*(-.*)*")

.PHONY: help
help: ## Show help
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo "${VERSION}"

.PHONY: test
test: ## Runs all tests
	@./test.sh

.PHONY: generate
generate: build ## Regenerates all files
	@PATH="$(PWD)/bin:$(PATH)" conflow generate

.PHONY: build
build: bin/conflow ## Build the conflow binary

bin/conflow:
	@echo "Building bin/conflow"
	@go version
	GOBIN="$(PWD)/bin" go install -ldflags="-s -w -X github.com/conflowio/conflow/src/conflow.Version=${VERSION}" ./cmd/conflow/

.PHONY: clean
clean: ## Clean all built files
	@rm -rf bin

.PHONY: clean-generated
clean-generated: ## Delete all generated files created by conflow
	@find . -name "*.cf.go" -type f -delete

.PHONY: goimports
goimports: ## Run goimports on all files
	@echo "Running goimports on all files"
	@./scripts/goimports.sh

.PHONY: lint
lint: ## Runs linting checks
	@echo "Running lint checks"
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

.PHONY: update-dependencies
update-dependencies: ## Updates all dependencies
	@echo "Updating Go dependencies"
	@cat go.mod | grep -E "^\t" | grep -v "// indirect" | cut -f 2 | cut -d ' ' -f 1 | xargs -n 1 -t go get -d -u
	@go mod vendor
	@go mod tidy

.PHONY: check
check: lint check-generate check-go-generate check-goimports ## Runs various code checks

.PHONY: check-generate
check-generate: generate
	@echo "Checking 'make generate'"
	@ scripts/check_git_changes.sh "make generate"

.PHONY: check-go-generate
check-go-generate:
	@echo "Checking 'go generate ./...'"
	@go generate ./...
	@ scripts/check_git_changes.sh "make go-generate"

.PHONY: check-goimports
check-goimports: goimports
	@echo "Checking 'make goimports"
	@ scripts/check_git_changes.sh "make goimports"
