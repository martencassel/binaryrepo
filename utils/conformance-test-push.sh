#!/bin/bash

export OCI_ROOT_URL=https://docker-local.example.com/
export OCI_TEST_PUSH=1
export OCI_BLOB_DIGEST="sha256:2834dc507516af02784808c5f48b7cbe38b8ed5d0f4837f16e78d00deb7e7767"
export OCI_BLOB_MANIFEST="sha256:2834dc507516af02784808c5f48b7cbe38b8ed5d0f4837f16e78d00deb7e7767"
export OCI_NAMESPACE="test"
export OCI_TEST_PULL=0
export OCI_TEST_CONTENT_DISCOVERY=0
export OCI_TEST_CONTENT_MANAGEMENT=0
export OCI_CROSSMOUNT_NAMESPACE="xyz"
export OCI_DEBUG=1
/tmp/distribution-spec/conformance/conformance.test

