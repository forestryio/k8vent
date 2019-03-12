# Makefile for Travis CI

GO = go
GO_FLAGS = -v
GO_ARGS = $(shell go list ./...)
GO_BUILD_ARGS =

TARGET = k8vent
DOCKER_TARGET = docker/$(TARGET)
DOCKER_IMAGE = atomist/$(TARGET)
DOCKER_VERSION = 0.11.0
DOCKER_TAG = $(DOCKER_IMAGE):$(DOCKER_VERSION)

all: vet

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

$(DOCKER_TARGET): clean-local
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GO_FLAGS) $(GO_BUILD_ARGS) -a --installsuffix cgo --ldflags="-s" -o "$(DOCKER_TARGET)"

docker-target: $(DOCKER_TARGET)

docker-build: docker-target
	cd docker && docker build -t "$(DOCKER_TAG)" .

docker: docker-build
	docker push "$(DOCKER_TAG)"

clean: clean-local
	$(GO) clean $(GO_FLAGS) $(GO_ARGS)

clean-local:
	-rm -f "$(DOCKER_TARGET)"

.PHONY: all fast clean build test vet
.PHONY: docker docker-target docker-build docker-push
.PHONY: clean-local
