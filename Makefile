
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p "$(LOCALBIN)"

SHELL := bash

CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint

# renovate: datasource=github-tags depName=kubernetes-sigs/controller-runtime
CONTROLLER_RUNTIME_VERSION = release-0.22

# renovate: datasource=github-tags depName=kubernetes-sigs/controller-tools extractVersion=^v\d+\.\d+\.\d+$
CONTROLLER_TOOLS_VERSION = v0.20.1

# renovate: datasource=github-tags depName=kubernetes-sigs/controller-tools extractVersion=^envtest-(?<version>v\d+\.\d+\.\d+)$
ENVTEST_K8S_VERSION = 1.35.0

# renovate: datasource=github-tags depName=golangci/golangci-lint
GOLANGCI_LINT_VERSION = v2.10.1

.DEFAULT_GOAL = help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


.PHONY: controller-gen
controller-gen:
	GOBIN="$(LOCALBIN)" go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: setup-envtest 
setup-envtest: 
	GOBIN="$(LOCALBIN)" go install sigs.k8s.io/controller-runtime/tools/setup-envtest@$(CONTROLLER_RUNTIME_VERSION)

.PHONY: golangci-lint 
golangci-lint:
	GOBIN="$(LOCALBIN)" go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: manifests
manifests: controller-gen
	"$(CONTROLLER_GEN)" rbac:roleName=ivu-manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen
	$(CONTROLLER_GEN) object paths=./...

.PHONY: test 
test: setup-envtest manifests generate ## Run the tests 
	KUBEBUILDER_ASSETS="$(shell setup-envtest use $(ENVTEST_K8S_VERSION) --bin-dir "$(LOCALBIN)" -p path)" go test -v ./...

.PHONY: test-short
test-short:  ## Skips slow integration tests
	go test -v ./... -short

.PHONY: clean
clean: setup-envtest ## Clean up envtest binaries
	$(ENVTEST) cleanup --bin-dir "$(LOCALBIN)"
	rm -rf $(LOCALBIN)

.PHONY: lint
lint: golangci-lint ## Run linter
	$(GOLANGCI_LINT) run

.PHONY: fmt
fmt: golangci-lint ## Run format
	$(GOLANGCI_LINT) fmt