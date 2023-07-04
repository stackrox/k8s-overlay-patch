$(call go-tool, CONTROLLER_GEN, sigs.k8s.io/controller-tools/cmd/controller-gen, tools)

.PHONY: download
download:
	@echo Download go.mod dependencies
	@go mod download

.PHONY: setup-tools
setup-tools: download
	cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: test
test:
	go test ./...

.PHONY: test-sanity
test-sanity: fix lint ## Test repo formatting, linting, etc.
	go vet ./...
	git diff --exit-code # diff again to ensure other checks don't change repo

.PHONY: fix
fix: setup-tools
	go mod tidy
	go fmt ./...
	golangci-lint run --fix

.PHONY: lint
lint: setup-tools
	golangci-lint run

.PHONY: all
all: setup-tools test lint

.PHONY: ci
ci: all test-sanity
