.PHONY: all build deps image lint migrate test vet
CHECK_FILES?=$$(go list ./... | grep -v /vendor/)

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: lint vet test build ## Run the tests and build the binary.

build: ## Build the binary.
	echo $$(git describe --tags --always --first-parent --dirty) > cmd/VERSION
	go build $(FLAGS)
	GOOS=linux GOARCH=arm64 go build $(FLAGS) -o supabase-admin-api-arm64

deps: ## Install dependencies.
	@go mod download

lint: ## Lint the code.
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.47.2 golangci-lint run -v --timeout 300s

test: ## Run tests.
	go test ./... -p 1 -race -v -count=1 -coverpkg ./cmd/...,./api/...,./optimizations/... -coverprofile=coverage.out

vet: # Vet the code
	go vet $(CHECK_FILES)
