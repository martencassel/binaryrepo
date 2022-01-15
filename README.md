# About this project

## Problem

Performance of package managers (Docker, Helm, Go etc) can be significantly improved by reusing previously fetched resources from the internet to a shared cache server.

Binaryrepo is a shared cache server that store responses to be reused by more than one user.
Multiple users may need to download a certain package from a repository on the internet.
By setting up a shared cache on the local network it may serve many users so that popular packages are
resused a number of times, reducing network traffic and latency.
## Current state

Binaryrepo can proxy Docker Hub, and it supports proxy caching images from docker hub.

Due to limitations in the docker client, a reverse proxy (nginx) must be setup infront of the binaryrepo server,
in order to be able to pull images through the binaryrepo server from docker hub.

The flow looks like this:

```bash
  docker image pull docker-remote.example.com ---> nginx:443 ---> binaryrepo:8081/repo/docker-remote/v2/*
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

## Run some tests
The following command will

1. Clear the local docker cache.
2. Pull images from binaryrepo (docker image pull) through nginx at docker-remote.example.com/v2/* endpoints
3. Binaryrepo will poulate it's cache in /tmp/filestore/*
4. The local docker cache will be cleared again.
5. Pull the same images as in step 2.
6. Now images will be served from binaryrepo's cache under /tmp/filestore/*

```bash
make test
```