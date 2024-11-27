LOCALBIN ?= $(shell pwd)/.bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool binary names.
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint
GO_FUMPT = $(LOCALBIN)/gofumpt
GCI = $(LOCALBIN)/gci

## Tool versions.
GOLANGCI_LINT_VERSION ?= v1.60.1
GO_FUMPT_VERSION ?= v0.6.0
GCI_VERSION ?= v0.13.5

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

.PHONY: gofumpt
gofumpt: $(GO_FUMPT)
$(GO_FUMPT): $(LOCALBIN)
	$(call go-install-tool,$(GO_FUMPT),mvdan.cc/gofumpt,$(GO_FUMPT_VERSION))

.PHONY: gci
gci: $(GCI)
$(GCI): $(LOCALBIN)
	$(call go-install-tool,$(GCI),github.com/daixiang0/gci,$(GCI_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef
