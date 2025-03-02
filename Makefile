VERSION := $(shell git describe --tags --always)

all: help

help: ## show help
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## clean artifacts
	@rm -rf ./coverage.txt ./*.out
	@rm -rf ./out ./.bin
	@echo Successfully removed artifacts

.PHONY: version
version: ## show version
	@echo $(VERSION)

.PHONY: lint
lint: ## run golangci-lint
	@golangci-lint run ./...

.PHONY: run-quest-1
run-quest-1:
	@go run quest1/*.go

.PHONY: run-quest-2
run-quest-2:
	@go run quest2/*.go

.PHONY: run-quest-3
run-quest-3:
	@go test -v quest3/*
