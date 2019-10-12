.DEFAULT_GOAL := help
SHELL := /usr/bin/env bash

.PHONY: help
help:
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Runs all tests
	@./test.sh

.PHONY: test
generate: bin/basil ## Regenerates all files
	@PATH="$(PWD)/bin;${PATH}" go generate ./...

.PHONY: build
build: bin/basil

bin/basil:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/basil ./cmd/basil/

.PHONY: clean
clean:
	rm -rf bin
