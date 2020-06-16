TARGET = k8svent
VERSION = $(shell git describe --always --dirty | sed 's/^v//')

GO = go
GO_FLAGS = -v
GO_ARGS = ./...
GO_BUILD_ARGS = -ldflags="-X github.com/atomist/k8svent/vent.Version=$(VERSION)"

GOLANGCI_LINT_ARGS = --timeout=2m

all: lint

generate:
	$(GO) generate $(GO_FLAGS) $(GO_ARGS)

build: generate
	$(GO) build $(GO_FLAGS) $(GO_BUILD_ARGS) -o "$(TARGET)"

test: build
	$(GO) test $(GO_FLAGS) $(GO_ARGS)

install: test
	$(GO) install $(GO_FLAGS) $(GO_ARGS)

vet: install
	$(GO) vet $(GO_FLAGS) $(GO_ARGS)

lint: vet
	golangci-lint run $(GOLANGCI_LINT_ARGS)

clean:
	$(GO) clean $(GO_FLAGS) $(GO_ARGS)

.PHONY: all build generate lint test vet
.PHONY: clean
