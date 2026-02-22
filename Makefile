
SHELL := bash
platform := $(shell uname | tr A-Z a-z)
ARCHITECTURE = $(shell go env GOARCH)

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p "$(LOCALBIN)"

ifeq ($(ARCHITECTURE),aarch64)
	ARCHITECTURE=arm64
endif

export KUBEBUILDER_ASSETS = $(LOCALBIN)/k8s/$(ENVTEST_K8S_VERSION)-$(platform)-$(ARCHITECTURE)
export PATH := $(LOCALBIN):$(PATH)


.DEFAULT_GOAL = help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# renovate: datasource=github-tags depName=kubernetes-sigs/controller-tools extractVersion=^envtest-(?<version>v\d+\.\d+\.\d+)$
ENVTEST_K8S_VERSION = 1.35.0

.PHONY: install-tools
install-tools: ## Install all tools
	cd internal/tools && cat tools.go | grep _ | awk -F'"' '{print $$2}' | GOBIN="$(LOCALBIN)" xargs -tI % go install -mod=mod %

.PHONY: manifests
manifests: install-tools
	controller-gen rbac:roleName=ivu-manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: install-tools
	controller-gen object paths=./...

.PHONY: test 
test: manifests generate ## Run the tests 
	 go test -v ./...

.PHONY: test-short
test-short:  ## Skips slow integration tests
	go test -v ./... -short

.PHONY: clean
clean: install-tools ## Clean up envtest binaries
	setup-envtest cleanup --bin-dir "$(LOCALBIN)"
	rm -rf $(LOCALBIN)

.PHONY: lint
lint: install-tools ## Run linter
	golangci-lint run

.PHONY: fmt
fmt: install-tools ## Run format
	golangci-lint fmt