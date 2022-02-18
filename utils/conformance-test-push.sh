#!/bin/bash
docker image pull docker-remote.example.com/alpine:latest

rm -rf /tmp/distribution-spec && cd /tmp && git clone https://github.com/opencontainers/distribution-spec.git && cd distribution-spec/conformance && go test -c

export OCI_ROOT_URL=http://localhost:8081/repo/docker-local
export OCI_TEST_PUSH=1
export OCI_BLOB_DIGEST="sha256:2834dc507516af02784808c5f48b7cbe38b8ed5d0f4837f16e78d00deb7e7767"
export OCI_BLOB_MANIFEST="sha256:2834dc507516af02784808c5f48b7cbe38b8ed5d0f4837f16e78d00deb7e7767"
export OCI_NAMESPACE="myrepo"
export OCI_TEST_PULL=0
export OCI_TEST_CONTENT_DISCOVERY=0
export OCI_TEST_CONTENT_MANAGEMENT=0
export OCI_CROSSMOUNT_NAMESPACE="docker-local"
export OCI_DEBUG=1
/tmp/distribution-spec/conformance/conformance.test

