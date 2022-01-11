#!/bin/bash

# Nginx need certificates. We create a local-ca in the home folder.

cd $HOME

docker run --rm -v $(pwd)/certs:/certs ehazlett/certm -d /certs \
    ca generate -o=local --overwrite

docker run --rm -v $(pwd)/certs:/certs ehazlett/certm -d /certs \
    server generate --host docker-remote.example.com -o=local --overwrite

# Fedora
sudo cp ./certs/ca.pem /etc/pki/ca-trust/source/anchors/
sudo update-ca-trust
