# About this project
## Motivation

Study various package managers and related protocols. For example, Docker Registry v2, Go Module registries, Helm registries etc.
Build a tool from scratch, simmiliar or identical to tools like Artifactory or Nexus (binary repository manager)
## Problem

Performance of package managers (docker, helm, or go and more) can be significantly improved by
resuing previously fetched resources from the internet to a local central shared cache server.

A shared cache stores responses to be reused by more than one user. For example, multiple users may need
to download the postgres:latest image from docker hub. By setting up a package cache server on the local network
it may serve many users so that popular docker images are reused a number of times, reducing
network traffic and latency.

This project tries to implement a binary artifact manager server that can serve software packages such as
Docker images, Helm charts or Go modules etc, either from a local storage or from remote sources (caching proxy).
All binary packages are stored in a "single-instance-store" using a checksum scheme.
## Current state

Binaryrepo is a server that serves a proxy cache for any docker registry that implements the Docker Registry v2 procotol.
In the current setup, public docker images can be pulled through it.

# Future plans

More features might be implemented.
# binary-repo

### Getting started

The following example will setup binaryrepo to be used
as a remote proxy cache of docker hub.

Docker pull command will access the remote repo through a local nginx container.

#### Prerequisites
1. Create certs
```bash
make setup-certs
ls ~/certs
```
2. Modify /etc/hosts
```bash
127.0.0.1 docker-remote.example.com
```
#### Start everything and run docker pull tests
```bash
make check-remote-pull
```
### Building

```bash
make build
```
### Working demo
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
