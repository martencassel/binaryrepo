go_version = 1.17.7

APP_VERSION?=$(shell git rev-parse --short HEAD)

GO_SRCS := $(shell find . -type f -name '*.go' -a ! -name 'zz_generated*')

# binaryrepo run on linux
GOARCH ?= $(shell go env GOARCH)
GOPATH ?= $(shell go env GOPATH)

GO ?= go
golint := $(shell which golangci-lint)

.PHONY: build
build: lint server
	docker-compose build

TARGET_OS ?= linux

binaryrepo: $(GO_SRCS)
			CGO_ENABLED=0 $(GO) build -o $(BUILD_DIR)/binaryrepo -ldflags "-X github.com/martencassel/binaryrepo/pkg/util/version.V=$(APP_VERSION)" ./cmd/binary-repo

.PHONY: build.docker-image
build.docker-image:
	docker build -t binaryrepo:latest -f ./docker/Dockerfile .

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
	bash ./scripts/docker-nginx/start-nginx.sh

.PHONY: start
start:  build
	unset http_proxy
	unset https_proxy
	./build/binaryrepo run > binaryrepo.log 2> binaryrepo.log &
	tail -f binaryrepo.log

.PHONY: stop
stop:
	kill `pidof binaryrepo`

.PHONY: setup-certs
setup-certs:
	bash ./scripts/docker-nginx/setup-certs.sh

.PHONY: up
up:
	docker-compose -f ./docker-compose.dev.yml up -d --remove-orphans

.PHONY: ps
ps:
	docker-compose -f ./docker-compose.dev.yml ps

.PHONY: down
down:
	docker-compose -f ./docker-compose.dev.yml down

.PHONY: db-shell
db-shell:
	docker exec -it binaryrepo-postgres psql -U postgres

.PHONY: logs
logs:
	docker-compose logs -f

.PHONY: clear-docker
clear-docker:
	docker rm -f $(docker ps -qa)||true


.PHONY: clear-local-images
clear-local-test-images:
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
server: clean
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/binaryrepo -ldflags "-X github.com/martencassel/binaryrepo/pkg/util/version.V=$(APP_VERSION)" ./cmd/binary-repo

.PHONY: clean
clean:
	rm -rf ./build

.PHONY: test
test: clear-local-test-images clear-filestore
	tree /tmp/filestore
	docker image pull docker-remote.example.com/redis:latest
	docker image pull docker-remote.example.com/postgres:latest
	tree /tmp/filestore
	docker rmi -f docker-remote.example.com/redis:latest
	docker rmi -f docker-remote.example.com/postgres:latest
	docker rmi -f redis:latest
	docker rmi -f postgres:latest
	docker image pull docker-remote.example.com/redis:latest
	docker image pull docker-remote.example.com/postgres:latest
