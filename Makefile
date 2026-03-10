
SHELL := bash
ENVTEST_K8S_VERSION := 1.35.0
KUBEBUILDER_PATH := $(shell setup-envtest use $(ENVTEST_K8S_VERSION) -i -p path)


.DEFAULT_GOAL = help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# renovate: datasource=github-tags depName=kubernetes-sigs/controller-tools extractVersion=^envtest-(?<version>v\d+\.\d+\.\d+)$

.PHONY: manifests
manifests:
	controller-gen rbac:roleName=ivu-manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate:
	controller-gen object paths=./...

.PHONY: test 
test: manifests generate ## Run the tests 
	KUBEBUILDER_ASSETS=$(KUBEBUILDER_PATH) go test -v ./...

.PHONY: test-short
test-short:  ## Skips slow integration tests
	go test -v ./... -short

.PHONY: clean
clean: ## Clean up envtest binaries
	setup-envtest cleanup

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: fmt
fmt: ## Run format
	golangci-lint fmt