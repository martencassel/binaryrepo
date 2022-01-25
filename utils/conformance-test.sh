#!/bin/bash

rm -rf /tmp/distribution-spec && cd /tmp && git clone https://github.com/opencontainers/distribution-spec.git && cd distribution-spec/conformance
go test -c
export OCI_ROOT_URL=http://localhost:8081/repo/docker-remote
export OCI_TEST_PUSH=0
export OCI_TEST_PULL=1
export OCI_TEST_CONTENT_DISCOVERY=0
export OCI_TEST_CONTENT_MANAGEMENT=0
export OCI_DEBUG=true
export OCI_NAMESPACE=alpine
./conformance.test

