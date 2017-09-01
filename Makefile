# Makefile for Travis CI

GO = go
GO_FLAGS = -v
GO_ARGS = +local
GO_BUILD_ARGS = +local,program
GOVENDOR = govendor

all: vet

fast:
	$(GO) build $(GO_FLAGS)

get: $(GOVENDOR)

generate: get
	$(GOVENDOR) generate $(GO_FLAGS) $(GO_ARGS)

build: generate
	$(GOVENDOR) build $(GO_FLAGS) $(GO_BUILD_ARGS)

test: build
	$(GOVENDOR) test $(GO_FLAGS) $(GO_ARGS)

vet: test
	$(GOVENDOR) vet $(GO_FLAGS) $(GO_ARGS)

docker: clean all
	docker login -u $$DOCKER_USER -p $$DOCKER_PASSWORD $$DOCKER_REGISTRY
	bash build-docker.bash

clean: $(GOVENDOR) clean-local
	$(GOVENDOR) clean $(GO_FLAGS) $(GO_ARGS)

clean-local:

$(GOVENDOR):
	$(GO) get $(GO_FLAGS) -u github.com/kardianos/govendor

.PHONY: all fast clean get build test vet docker
.PHONY: clean-local
.PHONY: govendor
