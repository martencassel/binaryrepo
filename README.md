# About this project
## Motivation

Study various package managers and related protocols. For example, Docker Registry v2, Go Module registries, Helm registries etc.
Build a tool from scratch, simmiliar or identical to tools like Artifactory or Nexus (binary repository manager)
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

## Issues

```bash
docker image pull docker-remote.example.com/postgres:latest
The layer 794976979956 times out, Its around 10*n MB.
Need to fix this:
https://developpaper.com/the-implementation-of-downloading-files-with-http-client-in-golang/
And check https://cs.opensource.google/go/go/+/go1.17.6:src/net/http/httputil/reverseproxy.go;l=143 implementation.
```

## Future plans
More features might be implemented.
## Getting started

The following example will setup binaryrepo to be used
as a remote proxy cache of docker hub.

Docker pull command will access the remote repo through a local nginx container.

## Prerequisites
1. Create certs
```bash
make setup-certs
ls ~/certs
```
2. Modify /etc/hosts
```bash
127.0.0.1 docker-remote.example.com
```
## Start everything and run docker pull tests
```bash
make check-remote-pull
```
## Building

```bash
make build
```
## Working demo
Proxy docker images from docker hub

- Create certs under $HOME/certs/ for nginx, add add the certs to your local cert store.
- Add docker-remote.example.com to your /etc/hosts file
- Start binaryrepo server at http://localhost:8081/
- Start nginx, it will proxies requests from docker-remote.example.com to http://localhost:8081/repo/docker-remote
- Run docker image pull on some sample images, remove the local images, and then pull again, now from the local cache in binaryrepo
  (tree /tmp/filestore)

```bash
make setup-certs
make check-remote-pull
```
