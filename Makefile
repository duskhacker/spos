FILTER := grep -v -e internal/ -e gen/ -e mock_
PKGS := $(shell go list ./... | $(FILTER))
VENDOR := govendor
GOLINT := golint
GOIMPORTS := goimports

# sources that don't match goimports
GOIMPORTS_SRCS := $(shell $(GOIMPORTS) -l . | $(FILTER))

# packages to test
TEST_PKGS := $(shell go list -f '{{.ImportPath}} {{.TestGoFiles}}' ./... | grep -v '\[\]$$' | cut -d' ' -f1 | grep -v -e /internal/ -e /gen/)

# flags passed to 'go test'
TEST_FLAGS := -short

GOLINT_FLAGS :=


default: goimports vendor vet lint test

test:
	-$(foreach p,${TEST_PKGS},go test $(p);)

vet:
	go vet $(PKGS)

lint: $(patsubst %,%.lint,$(PKGS))

%.lint:
	$(GOLINT) $(GOLINT_FLAGS) $* | grep -v 'should have comment' || exit 0

vendor:
	$(VENDOR) list | egrep -v "^i|^l|^p|^v" || exit 0

goimports: $(patsubst %,%.goimports,$(GOIMPORTS_SRCS))

%.goimports:
	@echo goimports -w $*

.PHONY: default generate test vet vendor goimports
