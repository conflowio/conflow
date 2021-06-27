.DEFAULT_GOAL := build

.PHONY: help
help:
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Runs all tests
	@./test.sh

.PHONY: generate
generate: build ## Regenerates all files
	@PATH="$(PWD)/bin:$(PATH)" basil generate

.PHONY: build
build: bin/basil ## Build the basil binary

bin/basil:
	@GOBIN="$(PWD)/bin" go install -ldflags="-s -w" ./cmd/basil/

.PHONY: clean
clean: ## Clean all built files
	@rm -rf bin

.PHONY: clean-generated
clean-generated: ## Delete all generated files created by basil
	@find . -name "*.basil.go" -type f -delete

.PHONY: goimports
goimports: ## Run goimports on all files
	@echo "Running goimports on all files"
	@./scripts/goimports.sh

.PHONY: lint
lint: ## Runs linting checks
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...
