#!/bin/bash
docker image pull docker-remote.example.com/alpine:latest

#rm -rf /tmp/distribution-spec && cd /tmp && git clone https://github.com/opencontainers/distribution-spec.git && cd distribution-spec/conformance && go test -c
#go test -c
export OCI_ROOT_URL=http://docker-remote.example.com:8081/repo/docker-remote
export OCI_TEST_PUSH=0
export OCI_TEST_PULL=1
export OCI_TEST_CONTENT_DISCOVERY=0
export OCI_TEST_CONTENT_MANAGEMENT=0
export OCI_DEBUG=true
export OCI_NAMESPACE=alpine
export OCI_BLOB_DIGEST="sha256:c059bfaa849c4d8e4aecaeb3a10c2d9b3d85f5165c66ad3a4d937758128c4d18"
export OCI_MANIFEST_DIGEST="sha256:e7d88de73db3d3fd9b2d63aa7f447a10fd0220b7cbf39803c803f2af9ba256b3"
export OCI_TAG_NAME="latest"
/tmp/distribution-spec/conformance/conformance.test

