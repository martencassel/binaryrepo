#!/bin/bash

rm -rf /tmp/distribution-spec && cd /tmp && git clone https://github.com/opencontainers/distribution-spec.git && cd distribution-spec/conformance && go test -c
