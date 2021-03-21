BINDIR  := $(CURDIR)/bin
BINNAME ?= uplift

GOBIN = $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN = $(shell go env GOPATH)/bin
endif

.PHONY: all
all: build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on go build -o '$(BINDIR)/$(BINNAME)' ./cmd/uplift

.PHONY: clean
clean:
	@rm -rf '$(BINDIR)'