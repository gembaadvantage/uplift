BINDIR     := $(CURDIR)/bin
BINNAME    ?= uplift
BINVERSION := ''
LDFLAGS    := -w -s

GOBIN = $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN = $(shell go env GOPATH)/bin
endif

# Interrogate git for build time information
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_BRANCH = $(shell git branch --show-current)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)

BINVERSION = ${GIT_TAG}

ifneq ($(GIT_BRANCH),'master')
	BINVERSION := $(BINVERSION)-${GIT_SHA}
endif

# Set build time information
LDFLAGS += -X github.com/gembaadvantage/uplift/internal/version.version=${BINVERSION}
LDFLAGS += -X github.com/gembaadvantage/uplift/internal/version.gitCommit=${GIT_COMMIT}
LDFLAGS += -X github.com/gembaadvantage/uplift/internal/version.gitBranch=${GIT_BRANCH}

.PHONY: all
all: build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on go build -ldflags '$(LDFLAGS)' -o '$(BINDIR)/$(BINNAME)' ./cmd/uplift

.PHONY: clean
clean:
	@rm -rf '$(BINDIR)'