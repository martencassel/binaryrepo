# About this project

## Problem

Performance of package managers (Docker, Helm, Go etc) can be significantly improved by reusing previously fetched resources from the internet to a shared cache server.

Binaryrepo is a shared cache server that store responses to be reused by more than one user.
Multiple users may need to download a certain package from a repository on the internet.
By setting up a shared cache on the local network it may serve many users so that popular packages are
resused a number of times, reducing network traffic and latency.
## Features

* Proxy caching images from docker hub using a docker hub login
* The local cache checksum based approach which reduces storage space by storing files only once.

## Demo

In the demo below, you can see that the the download time of postgres:latest will be reduced

[![asciicast](https://asciinema.org/a/1bHV8eIiAFO4t2G5Azx8HrqLs.svg)](https://asciinema.org/a/1bHV8eIiAFO4t2G5Azx8HrqLs)


## Reverse proxy and docker client

Due to limitations in the docker client, a reverse proxy (nginx) must be
used in order to proxy images through binaryrepo server.

The flow looks like this:

```bash
  1. docker image pull docker-remote.example.com sends requests to docker-remote.example.com that points to nginx.
     nginx listens on localhost:443.
  2. nginx sends request to binaryrepo at localhost:8081/repo/docker-remote/v2/*
  3. binaryrepo server authenticates to docker hub and forwards requests to docker hub
  4. binaryrepo servers receices image layers etc, and saves it to /tmp/filestore/ cache.
  5. the docker client is finally being served the pulled data.
```
## Getting started

The following example will setup binaryrepo to be used
as a remote proxy cache of docker hub.

Docker pull command will access the remote repo through a local nginx container.

## Prerequisites

Golang, git and make needs to be installed

Create a self-signed cert and add it to the trust store:
```bash
make setup-certs
ls ~/certs
# On fedora do this
sudo cp ./certs/ca.pem /etc/pki/ca-trust/source/anchors/
sudo update-ca-trust
```

Add a host entry for the reverse proxy:
```bash
127.0.0.1 docker-remote.example.com
```

Build the binary:
```bash
make build
```

Start nginx
```bash
make reverse-proxy
```

Currently binaryrepo only supports accessing docker hub using an hub account.
It's possible to pull images from docker hub, without an account. But this is not implemented yet.

These environment variables must be set before starting binaryrepo

```bash
export DOCKERHUB_USERNAME=<your username>
export DOCKERHUB_PASSWORD=<your password>
```

Start binaryrepo
```bash
make stop
make start
```
## Code

The docker proxy

https://github.com/martencassel/binaryrepo/tree/main/pkg/docker/proxy

The docker registry client

https://github.com/martencassel/binaryrepo/tree/main/pkg/docker/client

The filestore
https://github.com/martencassel/binaryrepo/tree/main/pkg/filestore/fs
