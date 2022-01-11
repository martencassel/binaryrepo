#!/bin/bash

# A certificate is issued to the pre-configured remote docker repository
# https://docker-remote.example.com.

# The following scripts creates a local ca authority in $HOME/certs folder.
# It also issues a server certificate to the nginx server that is issued to https://docker-remote.example.com

cd $HOME
sudo rm -rf $HOME/certs;

docker run --rm -v $(pwd)/certs:/certs ehazlett/certm -d /certs \
    ca generate -o=local --overwrite

docker run --rm -v $(pwd)/certs:/certs ehazlett/certm -d /certs \
    server generate --host docker-remote.example.com -o=local --overwrite

# Fedora
sudo cp ./certs/ca.pem /etc/pki/ca-trust/source/anchors/
sudo update-ca-trust

ls -l $HOME/certs
openssl x509 -in ~/certs/server.pem -text |grep example.com

