$(call go-tool, CONTROLLER_GEN, sigs.k8s.io/controller-tools/cmd/controller-gen, tools)

.PHONY: generate
generate:
	go run k8s.io/code-generator/cmd/deepcopy-gen@v0.27.3 -i github.com/stackrox/k8s-overlay-patch/pkg/apis/v1alpha1 -h hack/boilerplate.go.txt --trim-path-prefix github.com/stackrox/k8s-overlay-patch --output-base . --output-package github.com/stackrox/k8s-overlay-patch/pkg/apis/v1alpha1 -O zz_deepcopy

.PHONY: test
test:
	go test ./...

.PHONY: test-sanity
test-sanity: generate fix lint ## Test repo formatting, linting, etc.
	go vet ./...
	git diff --exit-code # diff again to ensure other checks don't change repo

.PHONY: fix
fix:
	go mod tidy
	go fmt ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3 run --fix

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3 run

.PHONY: all
all: generate test lint

