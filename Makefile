APP_VERSION?=$(shell git rev-parse --short HEAD)

GO_SRCS := $(shell find . -type f -name '*.go' -a ! -name 'zz_generated*')

# binaryrepo run on linux
GOARCH ?= $(shell go env GOARCH)
GOPATH ?= $(shell go env GOPATH)

GO ?= go
golint := $(shell which golangci-lint)

.PHONY: build
build: lint server


TARGET_OS ?= linux

.PHONY: check-remote-pull
check-remote-pull: build reverse-proxy start
	./inttest/remote_docker_pull.sh

.PHONY: lint
lint:
	$(golint) run --verbose ./...

.PHONY: check-unit
check-unit:
	go test ./... -v

.PHONY: cover
cover:
	go test ./... -cover

.PHONY: reverse-proxy
reverse-proxy:
	bash ./utils/docker-nginx/start-nginx.sh

.PHONY: start
start:  build
	unset http_proxy
	unset https_proxy
	./build//binaryrepo run > binaryrepo.log 2> binaryrepo.log &
	tail -f binaryrepo.log

.PHONY: stop
stop:
	kill `pidof binaryrepo`

.PHONY: setup-certs
setup-certs:
	bash ./utils/docker-nginx/setup-certs.sh

.PHONY: clear-local-images
clear-local-images:
	docker rmi -f docker-remote.example.com/redis:latest
	docker rmi -f docker-remote.example.com/postgres:latest
	docker rmi -f redis:latest
	docker rmi -f postgres:latest

.PHONY: clear-filestore
clear-filestore:
	rm -rf /tmp/filestore
	mkdir -p /tmp/filestore

BUILD_DIR := build

.PHONY: server
server:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/binaryrepo -ldflags "-X github.com/martencassel/binaryrepo/pkg/util/version.V=$(APP_VERSION)" ./cmd/binary-repo

.PHONY: clean
clean:
	rm -rf ./build

.PHONY: test
test:
	docker rmi -f docker-remote.example.com/redis:latest
	docker rmi -f docker-remote.example.com/postgres:latest
	docker image pull docker-remote.example.com/redis:latest
	docker image pull docker-remote.example.com/postgres:latest