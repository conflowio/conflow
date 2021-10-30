.DEFAULT_GOAL := help

.PHONY: help
help: ## Show help
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
	@echo "Building bin/basil"
	@go version
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
