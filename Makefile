$(call go-tool, CONTROLLER_GEN, sigs.k8s.io/controller-tools/cmd/controller-gen, tools)

.PHONY: download
download:
	@echo Download go.mod dependencies
	@go mod download

.PHONY: setup-tools
setup-tools: download
	cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: generate
generate: setup-tools
	deepcopy-gen -i github.com/stackrox/k8s-overlay-patch/pkg/apis/v1alpha1 -h hack/boilerplate.go.txt --trim-path-prefix github.com/stackrox/k8s-overlay-patch --output-base . --output-package github.com/stackrox/k8s-overlay-patch/pkg/apis/v1alpha1 -O zz_deepcopy

.PHONY: test
test:
	go test ./...

.PHONY: test-sanity
test-sanity: generate fix lint ## Test repo formatting, linting, etc.
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
all: generate test lint

