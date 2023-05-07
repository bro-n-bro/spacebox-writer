VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')

LD_FLAGS = -X github.com/hexy-dev/spacebox/spacebox-writer/version.Version=$(VERSION) \
	-X github.com/hexy-dev/spacebox/soacebox-writer/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

.PHONY: fix dep build test race lint stats

fix: ## Fix fieldalignment
	fieldalignment -fix ./...

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build the binary file
	@go build -v ./cmd/main.go

test: ## Run unittests
	@go test ./... -count=1 -coverprofile=coverage.out

make coverage: ## Run unittests with coverage
	@go tool cover -html=coverage.out -o coverage.html

race: dep ## Run data race detector
	@go test -race ./... -count=1

install-linter: ## Install golangci-lint
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.48.0

lint: install-linter ## Lint the files
	./scripts/golint.sh

stats: ## Code analytics
	scc .
