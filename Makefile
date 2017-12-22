# Makefile for Travis CI

GO = go
GO_FLAGS = -v
GO_ARGS =
GO_BUILD_ARGS =

all: vet

fast:
	$(GO) build $(GO_FLAGS)

generate:
	$(GO) generate $(GO_FLAGS) $(GO_ARGS) ./...

build: generate
	$(GO) build $(GO_FLAGS) $(GO_BUILD_ARGS)

test: build
	$(GO) test $(GO_FLAGS) $(GO_ARGS) ./...

vet: test
	$(GO) vet $(GO_FLAGS) $(GO_ARGS) ./...

docker: clean all
	@docker login -u $$DOCKER_USER -p $$DOCKER_PASSWORD $$DOCKER_REGISTRY
	bash build-docker.bash

clean: $(GO) clean-local
	$(GO) clean $(GO_FLAGS) $(GO_ARGS)

clean-local:

.PHONY: all fast clean build test vet docker
.PHONY: clean-local
