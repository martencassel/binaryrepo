#!/bin/bash

# Nginx need certificates. We create a local-ca in the home folder.

docker run --rm -v $(pwd)/certs:/certs ehazlett/certm -d /certs \
    ca generate -o=local --overwrite

docker run --rm -v $(pwd)/certs:/certs ehazlett/certm -d /certs \
    server generate --host project.example.com -o=local --overwrite

sudo cp ./certs/ca.pem /etc/pki/ca-trust/source/anchors/

sudo update-ca-trust
