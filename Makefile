GO_SRCS := $(shell find . -type f -name '*.go' -a ! -name 'zz_generated*')

# binaryrepo run on linux
GOARCH ?= $(shell go env GOARCH)
GOPATH ?= $(shell go env GOPATH)

GO ?= go
golint := $(shell which golangci-lint)

.PHONY: build
build: lint binaryrepo

TARGET_OS ?= linux

binaryrepo: $(GO_SRCS)
	GOOS=$(TARGET_OS) GOARCH=$(GOARCH) $(GO) build -o $@ main.go

.PHONY: check-remote-pull
check-remote-pull:
	./inttest/remote_docker_pull.sh

.PHONY: clean
clean:
	-rm -f ./binaryrepo

.PHONY: lint
lint:
	$(golint) run --verbose ./...

.PHONY: check-unit
check-unit:
	go test ./... -v

.PHONY: cover
cover:
	go test ./... -cover

.PHONY: run
run: binaryrepo
	unset http_proxy
	unset https_proxy
	./binaryrepo

.PHONY: setup-certs
setup-certs:
	bash ./utils/docker-nginx/setup-certs.sh

