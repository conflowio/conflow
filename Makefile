.DEFAULT_GOAL := help

.PHONY: help
help:
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Runs all tests
	@./test.sh

.PHONY: generate
generate: bin/basil ## Regenerates all files
	PATH="$(PWD)/bin:$(PATH)" go generate ./...

.PHONY: build
build: bin/basil

bin/basil:
	GOBIN="$(PWD)/bin" go install -ldflags="-s -w" ./cmd/basil/

.PHONY: clean
clean:
	rm -rf bin
